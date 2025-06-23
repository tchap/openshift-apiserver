package etcd

import (
	"github.com/openshift/openshift-apiserver/pkg/image/apiserver/internal/imageutil"
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	apirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistrytest "k8s.io/apiserver/pkg/registry/generic/testing"
	"k8s.io/apiserver/pkg/registry/rest"
	etcdtesting "k8s.io/apiserver/pkg/storage/etcd3/testing"
	"k8s.io/apiserver/pkg/storage/storagebackend"
	"k8s.io/kubernetes/pkg/api/legacyscheme"

	imageapi "github.com/openshift/openshift-apiserver/pkg/image/apis/image"
	"github.com/openshift/openshift-apiserver/pkg/image/apiserver/registry/image"

	// install all APIs
	_ "github.com/openshift/openshift-apiserver/pkg/api/install"
)

func newStorage(t *testing.T) (*REST, *etcdtesting.EtcdTestServer) {
	server, etcdStorage := etcdtesting.NewUnsecuredEtcd3TestClientServer(t)
	etcdStorage.Codec = legacyscheme.Codecs.LegacyCodec(schema.GroupVersion{Group: "image.openshift.io", Version: "v1"})
	etcdStorageConfigForImages := &storagebackend.ConfigForResource{Config: *etcdStorage, GroupResource: schema.GroupResource{Group: "image.openshift.io", Resource: "images"}}
	imageRESTOptions := generic.RESTOptions{StorageConfig: etcdStorageConfigForImages, Decorator: generic.UndecoratedStorage, DeleteCollectionWorkers: 1, ResourcePrefix: "images"}
	storage, err := NewREST(imageRESTOptions)
	if err != nil {
		t.Fatal(err)
	}
	return storage, server
}

func TestStorage(t *testing.T) {
	storage, _ := newStorage(t)
	image.NewRegistry(storage)
}

func TestCreate(t *testing.T) {
	storage, server := newStorage(t)
	defer server.Terminate(t)
	defer storage.Store.DestroyFunc()
	test := genericregistrytest.New(t, storage.Store).ClusterScope()
	valid := validImage()
	valid.Name = ""
	valid.GenerateName = "test-"
	test.TestCreate(
		valid,
		// invalid
		&imageapi.Image{},
	)
}

func TestUpdate(t *testing.T) {
	storage, server := newStorage(t)
	defer server.Terminate(t)
	defer storage.Store.DestroyFunc()
	test := genericregistrytest.New(t, storage.Store).ClusterScope()
	test.TestUpdate(
		validImage(),
		// updateFunc
		func(obj runtime.Object) runtime.Object {
			object := obj.(*imageapi.Image)
			object.DockerImageReference = "openshift/origin"
			return object
		},
		// invalid updateFunc
		func(obj runtime.Object) runtime.Object {
			object := obj.(*imageapi.Image)
			object.DockerImageReference = "\\"
			return object
		},
	)
}

func TestList(t *testing.T) {
	storage, server := newStorage(t)
	defer server.Terminate(t)
	defer storage.Store.DestroyFunc()
	test := genericregistrytest.New(t, storage.Store).ClusterScope()
	test.TestList(
		validImage(),
	)
}

func TestGet(t *testing.T) {
	storage, server := newStorage(t)
	defer server.Terminate(t)
	defer storage.Store.DestroyFunc()
	test := genericregistrytest.New(t, storage.Store).ClusterScope()
	test.TestGet(
		validImage(),
	)
}

func TestDelete(t *testing.T) {
	storage, server := newStorage(t)
	defer server.Terminate(t)
	defer storage.Store.DestroyFunc()
	test := genericregistrytest.New(t, storage.Store).ClusterScope()
	image := validImage()
	image.ObjectMeta = metav1.ObjectMeta{GenerateName: "foo"}
	test.TestDelete(
		validImage(),
	)
}

func TestWatch(t *testing.T) {
	storage, server := newStorage(t)
	defer server.Terminate(t)
	defer storage.Store.DestroyFunc()
	test := genericregistrytest.New(t, storage.Store)

	valid := validImage()
	valid.Name = "foo"
	valid.Labels = map[string]string{"foo": "bar"}

	test.TestWatch(
		valid,
		// matching labels
		[]labels.Set{{"foo": "bar"}},
		// not matching labels
		[]labels.Set{{"foo": "baz"}},
		// matching fields
		[]fields.Set{
			{"metadata.name": "foo"},
		},
		// not matchin fields
		[]fields.Set{
			{"metadata.name": "bar"},
		},
	)
}

func TestCreateSetsMetadata(t *testing.T) {
	testCases := []struct {
		name   string
		image  *imageapi.Image
		expect func(*imageapi.Image) bool
	}{
		{
			name: "image schema v2",
			expect: func(image *imageapi.Image) bool {
				ok := true
				if image.DockerImageMetadata.Size != 67218142 {
					t.Errorf("image had size %d", image.DockerImageMetadata.Size)
					ok = false
				}
				if len(image.DockerImageLayers) != 1 || image.DockerImageLayers[0].Name != "sha256:a425929bf30ed86191ff0d58db6907a92e41de14e58465658e2c165b391cd49d" || image.DockerImageLayers[0].LayerSize != 67215871 {
					t.Errorf("unexpected layers: %#v", image.DockerImageLayers)
					ok = false
				}
				return ok
			},
			image: etcdImage(),
		},
	}

	for i, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			storage, server := newStorage(t)
			defer server.Terminate(t)
			defer storage.Store.DestroyFunc()

			obj, err := storage.Create(apirequest.NewDefaultContext(), test.image, rest.ValidateAllObjectFunc, &metav1.CreateOptions{})
			if obj == nil {
				t.Errorf("%d: Expected nil obj, got %v", i, obj)
				return
			}
			if err != nil {
				t.Errorf("%d: Unexpected non-nil error: %#v", i, err)
				return
			}
			image, ok := obj.(*imageapi.Image)
			if !ok {
				t.Errorf("%d: Expected image type, got: %#v", i, obj)
				return
			}
			if test.expect != nil && !test.expect(image) {
				t.Errorf("%d: Unexpected image: %#v", i, obj)
			}
		})
	}
}

func TestUpdateResetsMetadata(t *testing.T) {
	testCases := []struct {
		name     string
		image    *imageapi.Image
		existing *imageapi.Image
		expect   func(*imageapi.Image) bool
	}{
		{
			name: "manifest changes are ignored",
			expect: func(image *imageapi.Image) bool {
				ok := true
				if image.Labels["a"] != "b" {
					t.Errorf("unexpected labels: %s", image.Labels)
					ok = false
				}
				if image.DockerImageMetadata.ID != etcdConfigDigest {
					t.Errorf("unexpected container image: %#v", image.DockerImageMetadata)
					ok = false
				}
				if image.DockerImageReference != "bitnami/etcd-updated" {
					t.Errorf("image reference not changed: %s", image.DockerImageReference)
					ok = false
				}
				if image.DockerImageMetadata.Size != 67218142 {
					t.Errorf("image had size %d", image.DockerImageMetadata.Size)
					ok = false
				}
				if len(image.DockerImageLayers) != 1 || image.DockerImageLayers[0].LayerSize != 67215871 {
					t.Errorf("unexpected layers: %#v", image.DockerImageLayers)
					ok = false
				}
				return ok
			},
			existing: etcdImage(),
			image: etcdImage(func(img *imageapi.Image) {
				img.Labels = map[string]string{"a": "b"}
				img.DockerImageReference = "bitnami/etcd-updated"
			}),
		},
		{
			name: "manifest is preserved and unpacked",
			expect: func(image *imageapi.Image) bool {
				if len(image.DockerImageManifest) != 0 {
					t.Errorf("unexpected not empty manifest")
					return false
				}
				if image.DockerImageMetadata.ID != etcdConfigDigest {
					t.Errorf("unexpected container image: %#v", image.DockerImageMetadata)
					return false
				}
				if image.DockerImageReference != "bitnami/etcd-updated" {
					t.Errorf("image reference not changed: %s", image.DockerImageReference)
					return false
				}
				if image.DockerImageMetadata.Size != 67218142 {
					t.Errorf("image had size %d", image.DockerImageMetadata.Size)
					return false
				}
				if len(image.DockerImageLayers) != 1 || image.DockerImageLayers[0].Name != "sha256:a425929bf30ed86191ff0d58db6907a92e41de14e58465658e2c165b391cd49d" || image.DockerImageLayers[0].LayerSize != 67215871 {
					t.Errorf("unexpected layers: %#v", image.DockerImageLayers)
					return false
				}
				return true
			},
			existing: etcdImage(),
			image: &imageapi.Image{
				ObjectMeta: metav1.ObjectMeta{
					Name:         etcdManifestDigest,
					GenerateName: "etcd",
				},
				DockerImageReference: "bitnami/etcd-updated",
			},
		},
	}

	for i, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			storage, server := newStorage(t)
			defer server.Terminate(t)
			defer storage.Store.DestroyFunc()

			// Clear the resource version before creating
			test.existing.ResourceVersion = ""
			created, err := storage.Create(apirequest.NewDefaultContext(), test.existing, rest.ValidateAllObjectFunc, &metav1.CreateOptions{})
			if err != nil {
				t.Errorf("%d: Unexpected non-nil error: %#v", i, err)
				return
			}

			// Copy the resource version into our update object
			test.image.ResourceVersion = created.(*imageapi.Image).ResourceVersion
			obj, _, err := storage.Update(apirequest.NewDefaultContext(), test.image.Name, rest.DefaultUpdatedObjectInfo(test.image), rest.ValidateAllObjectFunc, rest.ValidateAllObjectUpdateFunc, false, &metav1.UpdateOptions{})
			if err != nil {
				t.Errorf("%d: Unexpected non-nil error: %#v", i, err)
				return
			}
			if obj == nil {
				t.Errorf("%d: Expected nil obj, got %v", i, obj)
				return
			}
			image, ok := obj.(*imageapi.Image)
			if !ok {
				t.Errorf("%d: Expected image type, got: %#v", i, obj)
				return
			}
			if test.expect != nil && !test.expect(image) {
				t.Errorf("%d: Unexpected image: %#v", i, obj)
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

const etcdManifestDigest = `sha256:c5f810e973f9d0985c4075650cbfaf6c2fb31d6508fdb2e5ed03dc58c5f5d2a6`

const etcdManifest = `{
    "schemaVersion": 2,
    "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
    "config": {
        "mediaType": "application/vnd.docker.container.image.v1+json",
        "size": 7429,
        "digest": "sha256:f21cd1b754d0525496ccae928c01718a5cc7596f5b19ddd8c8f522aee91adff2"
    },
    "layers": [
        {
            "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
            "size": 67215871,
            "digest": "sha256:a425929bf30ed86191ff0d58db6907a92e41de14e58465658e2c165b391cd49d"
        }
    ]
}`

const etcdConfigDigest = `sha256:f21cd1b754d0525496ccae928c01718a5cc7596f5b19ddd8c8f522aee91adff2`

const etcdConfig = `{
    "architecture": "amd64",
    "created": "2025-06-13T11:40:15.124280024Z",
    "history": [
        {
            "created": "0001-01-01T00:00:00Z",
            "created_by": "crane flatten sha256:9795186cd4d9720cb833e3ee623fafa358060e7567fbb9480cb43d75dbac6f8b"
        }
    ],
    "os": "linux",
    "rootfs": {
        "type": "layers",
        "diff_ids": [
            "sha256:c05fe3cfa98ad11079103bd5d9b6d784aa3030d9a273e8fc5cfbeaa813b88d69"
        ]
    },
    "config": {
        "Cmd": [
            "/opt/bitnami/scripts/etcd/run.sh"
        ],
        "Entrypoint": [
            "/opt/bitnami/scripts/etcd/entrypoint.sh"
        ],
        "Env": [
            "PATH=/opt/bitnami/common/bin:/opt/bitnami/etcd/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
            "HOME=/",
            "OS_ARCH=amd64",
            "OS_FLAVOUR=debian-12",
            "OS_NAME=linux",
            "APP_VERSION=3.6.1",
            "BITNAMI_APP_NAME=etcd"
        ],
        "Labels": {
            "com.vmware.cp.artifact.flavor": "sha256:c50c90cfd9d12b445b011e6ad529f1ad3daea45c26d20b00732fae3cd71f6a83",
            "org.opencontainers.image.base.name": "docker.io/bitnami/minideb:bookworm",
            "org.opencontainers.image.created": "2025-06-13T11:21:11Z",
            "org.opencontainers.image.description": "Application packaged by Broadcom, Inc.",
            "org.opencontainers.image.documentation": "https://github.com/bitnami/containers/tree/main/bitnami/etcd/README.md",
            "org.opencontainers.image.ref.name": "3.6.1-debian-12-r2",
            "org.opencontainers.image.source": "https://github.com/bitnami/containers/tree/main/bitnami/etcd",
            "org.opencontainers.image.title": "etcd",
            "org.opencontainers.image.vendor": "Broadcom, Inc.",
            "org.opencontainers.image.version": "3.6.1"
        },
        "User": "1001",
        "WorkingDir": "/opt/bitnami/etcd",
        "ExposedPorts": {
            "2379/tcp": {},
            "2380/tcp": {}
        },
        "ArgsEscaped": true,
        "Shell": [
            "/bin/bash",
            "-o",
            "errexit",
            "-o",
            "nounset",
            "-o",
            "pipefail",
            "-c"
        ]
    }
}`

var etcdImageTemplate = imageapi.Image{
	ObjectMeta: metav1.ObjectMeta{
		Name:         etcdManifestDigest,
		GenerateName: "etcd",
	},
	DockerImageReference:         "bitnami/etcd",
	DockerImageManifestMediaType: "application/vnd.docker.distribution.manifest.v2+json",
	DockerImageManifest:          etcdManifest,
	DockerImageConfig:            etcdConfig,
}

func etcdImage(hooks ...func(*imageapi.Image)) *imageapi.Image {
	img := etcdImageTemplate.DeepCopy()
	for _, h := range hooks {
		h(img)
	}
	if err := imageutil.InternalImageWithMetadata(img); err != nil {
		panic(err)
	}
	return img
}
