package imagestreamimage

import (
	"context"
	"testing"
	"time"

	etcd "go.etcd.io/etcd/client/v3"
	authorizationapi "k8s.io/api/authorization/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/diff"
	apirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/generic"
	etcdtesting "k8s.io/apiserver/pkg/storage/etcd3/testing"
	"k8s.io/apiserver/pkg/storage/storagebackend"
	"k8s.io/kubernetes/pkg/api/legacyscheme"

	imagev1 "github.com/openshift/api/image/v1"
	imageapi "github.com/openshift/openshift-apiserver/pkg/image/apis/image"
	"github.com/openshift/openshift-apiserver/pkg/image/apis/image/validation/fake"
	admfake "github.com/openshift/openshift-apiserver/pkg/image/apiserver/admission/fake"
	"github.com/openshift/openshift-apiserver/pkg/image/apiserver/registry/image"
	imageetcd "github.com/openshift/openshift-apiserver/pkg/image/apiserver/registry/image/etcd"
	"github.com/openshift/openshift-apiserver/pkg/image/apiserver/registry/imagestream"
	imagestreametcd "github.com/openshift/openshift-apiserver/pkg/image/apiserver/registry/imagestream/etcd"
	"github.com/openshift/openshift-apiserver/pkg/image/apiserver/registryhostname"

	_ "github.com/openshift/openshift-apiserver/pkg/api/install"
)

type fakeSubjectAccessReviewRegistry struct{}

func (f *fakeSubjectAccessReviewRegistry) Create(_ context.Context, subjectAccessReview *authorizationapi.SubjectAccessReview, _ metav1.CreateOptions) (*authorizationapi.SubjectAccessReview, error) {
	return nil, nil
}

func (f *fakeSubjectAccessReviewRegistry) CreateContext(ctx context.Context, subjectAccessReview *authorizationapi.SubjectAccessReview) (*authorizationapi.SubjectAccessReview, error) {
	return f.Create(ctx, subjectAccessReview, metav1.CreateOptions{})
}

func TestGet(t *testing.T) {
	tests := map[string]struct {
		input                   string
		repo                    *imageapi.ImageStream
		images                  []*imageapi.Image
		expectedImageMetadataID string
		expectedConfigHostname  string
		expectError             bool
	}{
		"empty string": {
			input:       "",
			expectError: true,
		},
		"one part": {
			input:       "a",
			expectError: true,
		},
		"more than 2 parts": {
			input:       "a@b@c",
			expectError: true,
		},
		"empty name part": {
			input:       "@id",
			expectError: true,
		},
		"empty id part": {
			input:       "name@",
			expectError: true,
		},
		"repo not found": {
			input:       "repo@id",
			repo:        nil,
			expectError: true,
		},
		"nil tags": {
			input:       "repo@id",
			repo:        &imageapi.ImageStream{},
			expectError: true,
		},
		"image not found": {
			input: "repo@id",
			repo: &imageapi.ImageStream{
				Status: imageapi.ImageStreamStatus{
					Tags: map[string]imageapi.TagEventList{
						"latest": {
							Items: []imageapi.TagEvent{
								{Image: "anotherid"},
							},
						},
					},
				},
			},
			expectError: true,
		},
		"happy path": {
			input: "repo@" + ubuntuManifestDigest,
			repo: &imageapi.ImageStream{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "ns",
					Name:      "repo",
				},
				Status: imageapi.ImageStreamStatus{
					Tags: map[string]imageapi.TagEventList{
						"latest": {
							Items: []imageapi.TagEvent{
								{Image: "anotherid"},
								{Image: "anotherid2"},
								{Image: ubuntuManifestDigest},
							},
						},
					},
				},
			},
			expectedImageMetadataID: ubuntuConfigDigest,
			expectedConfigHostname:  "57d2303e19c8",
			images: []*imageapi.Image{
				validImage(),
			},
		},
		"uses annotations from image stream": {
			input:       "repo@sha256:abc321",
			expectError: false,
			repo: &imageapi.ImageStream{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "ns",
					Name:      "repo",
					Annotations: map[string]string{
						"test":         "123",
						"another-test": "abc",
					},
				},
				Status: imageapi.ImageStreamStatus{
					Tags: map[string]imageapi.TagEventList{
						"latest": {
							Items: []imageapi.TagEvent{
								{Image: "anotherid"},
								{Image: "anotherid2"},
								{Image: "sha256:abc321"},
							},
						},
					},
				},
			},
			images: []*imageapi.Image{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "sha256:abc321",
					},
				},
			},
		},
		"matches partial sha": {
			input:       "repo@sha256:ff46b782",
			expectError: false,
			repo: &imageapi.ImageStream{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "ns",
					Name:      "repo",
				},
				Status: imageapi.ImageStreamStatus{
					Tags: map[string]imageapi.TagEventList{
						"latest": {
							Items: []imageapi.TagEvent{
								{Image: "anotherid"},
								{Image: "anotherid2"},
								{Image: "sha256:ff46b78279f207db3b8e57e20dee7cecef3567d09489369d80591f150f9c8154"},
							},
						},
					},
				},
			},
			images: []*imageapi.Image{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "sha256:ff46b78279f207db3b8e57e20dee7cecef3567d09489369d80591f150f9c8154",
					},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			server, etcdStorage := etcdtesting.NewUnsecuredEtcd3TestClientServer(t)
			defer server.Terminate(t)
			etcdStorage.Codec = legacyscheme.Codecs.LegacyCodec(
				schema.GroupVersion{Group: "image.openshift.io", Version: "v1"})
			client := etcd.NewKV(server.V3Client.Client)
			etcdStorageConfigForImages := &storagebackend.ConfigForResource{
				Config:        *etcdStorage,
				GroupResource: schema.GroupResource{Group: "image.openshift.io", Resource: "images"},
			}
			imageRESTOptions := generic.RESTOptions{
				StorageConfig:           etcdStorageConfigForImages,
				Decorator:               generic.UndecoratedStorage,
				DeleteCollectionWorkers: 1,
				ResourcePrefix:          "images",
			}
			imageStorage, err := imageetcd.NewREST(imageRESTOptions)
			if err != nil {
				t.Fatal(err)
			}
			defaultRegistry := registryhostname.DefaultRegistryHostnameRetriever("", "defaultregistry:5000")
			etcdStorageConfigForImageStreams := &storagebackend.ConfigForResource{
				Config:        *etcdStorage,
				GroupResource: schema.GroupResource{Group: "image.openshift.io", Resource: "imagestreams"},
			}
			imagestreamRESTOptions := generic.RESTOptions{
				StorageConfig:           etcdStorageConfigForImageStreams,
				Decorator:               generic.UndecoratedStorage,
				DeleteCollectionWorkers: 1,
				ResourcePrefix:          "imagestreams",
			}
			imageIndex := imagestreametcd.NewMockImageLayerIndex()
			imageStreamStorage, imageStreamLayersStorage, imageStreamStatus, internalStorage, err := imagestreametcd.NewRESTWithLimitVerifier(
				imagestreamRESTOptions,
				defaultRegistry,
				&fakeSubjectAccessReviewRegistry{},
				&admfake.ImageStreamLimitVerifier{},
				&fake.RegistryWhitelister{},
				imageIndex,
			)
			if err != nil {
				t.Fatal(err)
			}

			imageRegistry := image.NewRegistry(imageStorage)
			imageStreamRegistry := imagestream.NewRegistry(
				imageStreamStorage,
				imageStreamStatus,
				internalStorage,
				imageStreamLayersStorage,
			)

			storage := NewREST(imageRegistry, imageStreamRegistry)
			ctx := apirequest.NewDefaultContext()

			if test.repo != nil {
				ctx = apirequest.WithNamespace(apirequest.NewContext(), test.repo.Namespace)
				_, err := client.Put(
					context.TODO(),
					etcdtesting.AddPrefix("/imagestreams/"+test.repo.Namespace+"/"+test.repo.Name),
					runtime.EncodeOrDie(legacyscheme.Codecs.LegacyCodec(imagev1.SchemeGroupVersion), test.repo),
				)
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
					return
				}
			}
			if len(test.images) > 0 {
				for _, image := range test.images {
					_, err := client.Put(
						context.TODO(),
						etcdtesting.AddPrefix("/images/"+image.Name),
						runtime.EncodeOrDie(
							legacyscheme.Codecs.LegacyCodec(imagev1.SchemeGroupVersion),
							image,
						),
					)
					if err != nil {
						t.Fatalf("Unexpected error: %v", err)
					}
					imageIndex.Add(&imagev1.Image{
						ObjectMeta:          image.ObjectMeta,
						DockerImageManifest: image.DockerImageManifest,
					})
				}
			}

			obj, err := storage.Get(ctx, test.input, &metav1.GetOptions{})
			if test.expectError {
				if err == nil {
					t.Fatal("expected error but didn't get one")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %#v", err)
			}

			imageStreamImage := obj.(*imageapi.ImageStreamImage)
			// validate a couple of the fields
			if e, a := test.repo.Namespace, "ns"; e != a {
				t.Errorf("%s: namespace: expected %q, got %q", name, e, a)
			}
			if e, a := test.input, imageStreamImage.Name; e != a {
				t.Errorf("%s: name: expected %q, got %q", name, e, a)
			}

			expectedAnnotations := test.repo.ObjectMeta.Annotations
			gotAnnotations := imageStreamImage.ObjectMeta.Annotations
			if !equality.Semantic.DeepEqual(expectedAnnotations, gotAnnotations) {
				t.Error("Expected image stream annotations to match image stream image's")
				t.Log(diff.ObjectGoPrintDiff(expectedAnnotations, gotAnnotations))
			}

			expectedID := test.expectedImageMetadataID
			if expectedID != "" && expectedID != imageStreamImage.Image.DockerImageMetadata.ID {
				t.Errorf("id: expected %q, got %q", expectedID, imageStreamImage.Image.DockerImageMetadata.ID)
			}
			expectedConfigHostname := test.expectedConfigHostname
			hostname := imageStreamImage.Image.DockerImageMetadata.ContainerConfig.Hostname
			if expectedConfigHostname != "" && expectedConfigHostname != hostname {
				t.Errorf("container config hostname: expected %q, got %q", expectedConfigHostname, hostname)
			}
		})
	}
}

const ubuntuManifestDigest = `sha256:04f510bf1f2528604dc2ff46b517dbdbb85c262d62eacc4aa4d3629783036096`

const ubuntuManifest = `{
    "schemaVersion": 2,
    "mediaType": "application/vnd.oci.image.manifest.v1+json",
    "config": {
        "mediaType": "application/vnd.oci.image.config.v1+json",
        "size": 2295,
        "digest": "sha256:bf16bdcff9c96b76a6d417bd8f0a3abe0e55c0ed9bdb3549e906834e2592fd5f"
    },
    "layers": [
        {
            "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
            "size": 29715337,
            "digest": "sha256:d9d352c11bbd3880007953ed6eec1cbace76898828f3434984a0ca60672fdf5a"
        }
    ]
}`

const ubuntuConfigDigest = `sha256:bf16bdcff9c96b76a6d417bd8f0a3abe0e55c0ed9bdb3549e906834e2592fd5f`

const ubuntuConfig = `{
    "architecture": "amd64",
    "config": {
        "Hostname": "",
        "Domainname": "",
        "User": "",
        "AttachStdin": false,
        "AttachStdout": false,
        "AttachStderr": false,
        "Tty": false,
        "OpenStdin": false,
        "StdinOnce": false,
        "Env": [
            "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
        ],
        "Cmd": [
            "/bin/bash"
        ],
        "Image": "sha256:825befda5d2b1a76b71f4e1d6d31f5d82d4488b8337b1ad42e29b1340d766647",
        "Volumes": null,
        "WorkingDir": "",
        "Entrypoint": null,
        "OnBuild": null,
        "Labels": {
            "org.opencontainers.image.ref.name": "ubuntu",
            "org.opencontainers.image.version": "24.04"
        }
    },
    "container": "57d2303e19c80641e487894fdb01e8e26ab42726f45e72624efe9d812e1c8889",
    "container_config": {
        "Hostname": "57d2303e19c8",
        "Domainname": "",
        "User": "",
        "AttachStdin": false,
        "AttachStdout": false,
        "AttachStderr": false,
        "Tty": false,
        "OpenStdin": false,
        "StdinOnce": false,
        "Env": [
            "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
        ],
        "Cmd": [
            "/bin/sh",
            "-c",
            "#(nop) ",
            "CMD [\"/bin/bash\"]"
        ],
        "Image": "sha256:825befda5d2b1a76b71f4e1d6d31f5d82d4488b8337b1ad42e29b1340d766647",
        "Volumes": null,
        "WorkingDir": "",
        "Entrypoint": null,
        "OnBuild": null,
        "Labels": {
            "org.opencontainers.image.ref.name": "ubuntu",
            "org.opencontainers.image.version": "24.04"
        }
    },
    "created": "2025-05-29T04:21:01.971275965Z",
    "docker_version": "24.0.7",
    "history": [
        {
            "created": "2025-05-29T04:20:59.390476489Z",
            "created_by": "/bin/sh -c #(nop)  ARG RELEASE",
            "empty_layer": true
        },
        {
            "created": "2025-05-29T04:20:59.425928067Z",
            "created_by": "/bin/sh -c #(nop)  ARG LAUNCHPAD_BUILD_ARCH",
            "empty_layer": true
        },
        {
            "created": "2025-05-29T04:20:59.461048974Z",
            "created_by": "/bin/sh -c #(nop)  LABEL org.opencontainers.image.ref.name=ubuntu",
            "empty_layer": true
        },
        {
            "created": "2025-05-29T04:20:59.498669132Z",
            "created_by": "/bin/sh -c #(nop)  LABEL org.opencontainers.image.version=24.04",
            "empty_layer": true
        },
        {
            "created": "2025-05-29T04:21:01.6549815Z",
            "created_by": "/bin/sh -c #(nop) ADD file:598ca0108009b5c2e9e6f4fc4bd19a6bcd604fccb5b9376fac14a75522a5cfa3 in / "
        },
        {
            "created": "2025-05-29T04:21:01.971275965Z",
            "created_by": "/bin/sh -c #(nop)  CMD [\"/bin/bash\"]",
            "empty_layer": true
        }
    ],
    "os": "linux",
    "rootfs": {
        "type": "layers",
        "diff_ids": [
            "sha256:a8346d259389bc6221b4f3c61bad4e48087c5b82308e8f53ce703cfc8333c7b3"
        ]
    }
}`

func validImage() *imageapi.Image {
	return &imageapi.Image{
		ObjectMeta: metav1.ObjectMeta{
			Name:         ubuntuManifestDigest,
			GenerateName: "foo",
		},
		DockerImageReference:         "openshift/origin",
		DockerImageManifestMediaType: "application/vnd.oci.image.manifest.v1+json",
		DockerImageManifest:          ubuntuManifest,
		DockerImageConfig:            ubuntuConfig,
		DockerImageLayers: []imageapi.ImageLayer{
			{
				Name:      "sha256:d9d352c11bbd3880007953ed6eec1cbace76898828f3434984a0ca60672fdf5a",
				LayerSize: 29715337,
				MediaType: "application/vnd.oci.image.layer.v1.tar+gzip",
			},
		},
		DockerImageMetadata: imageapi.DockerImage{
			ID:            ubuntuConfigDigest,
			Parent:        "",
			Comment:       "",
			Created:       metav1.Date(2025, 5, 29, 4, 21, 1, 971275965, time.UTC),
			Container:     "57d2303e19c80641e487894fdb01e8e26ab42726f45e72624efe9d812e1c8889",
			DockerVersion: "24.0.7",
			Author:        "",
			Architecture:  "amd64",
			Size:          29718508,
			ContainerConfig: imageapi.DockerConfig{
				Hostname:        "57d2303e19c8",
				Domainname:      "",
				User:            "",
				Memory:          0,
				MemorySwap:      0,
				CPUShares:       0,
				CPUSet:          "",
				AttachStdin:     false,
				AttachStdout:    false,
				AttachStderr:    false,
				PortSpecs:       nil,
				ExposedPorts:    nil,
				Tty:             false,
				OpenStdin:       false,
				StdinOnce:       false,
				Env:             []string{"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"},
				Cmd:             []string{"/bin/sh", "-c", "#(nop) ", `CMD ["/bin/bash"]`},
				Image:           "sha256:825befda5d2b1a76b71f4e1d6d31f5d82d4488b8337b1ad42e29b1340d766647",
				Volumes:         nil,
				WorkingDir:      "",
				Entrypoint:      nil,
				NetworkDisabled: false,
				SecurityOpts:    nil,
				OnBuild:         nil,
				Labels: map[string]string{
					"org.opencontainers.image.ref.name": "ubuntu",
					"org.opencontainers.image.version":  "24.04",
				},
			},
			Config: &imageapi.DockerConfig{
				Hostname:        "",
				Domainname:      "",
				User:            "",
				Memory:          0,
				MemorySwap:      0,
				CPUShares:       0,
				CPUSet:          "",
				AttachStdin:     false,
				AttachStdout:    false,
				AttachStderr:    false,
				PortSpecs:       nil,
				ExposedPorts:    nil,
				Tty:             false,
				OpenStdin:       false,
				StdinOnce:       false,
				Env:             []string{"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"},
				Cmd:             []string{"/bin/bash"},
				Image:           "sha256:825befda5d2b1a76b71f4e1d6d31f5d82d4488b8337b1ad42e29b1340d766647",
				Volumes:         nil,
				WorkingDir:      "",
				Entrypoint:      nil,
				NetworkDisabled: false,
				OnBuild:         nil,
				Labels: map[string]string{
					"org.opencontainers.image.ref.name": "ubuntu",
					"org.opencontainers.image.version":  "24.04",
				},
			},
		},
	}
}
