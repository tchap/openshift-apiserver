package testutil

import (
	imageapi "github.com/openshift/openshift-apiserver/pkg/image/apis/image"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

// kindest/node:v1.33.1
const ImageSchemaV2ManifestDigest = `sha256:14ffd6ee8a3daa20cc934ba786626b181e1797268c5465f2c299a7cf54494c77`

const kindestManifest = `{
    "schemaVersion": 2,
    "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
    "config": {
        "mediaType": "application/vnd.docker.container.image.v1+json",
        "size": 1984,
        "digest": "sha256:d6b20550c77b11385dd30115ba29dbf9a9bfc98c2f28ff7d162a6ad7c9686251"
    },
    "layers": [
        {
            "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
            "size": 132852852,
            "digest": "sha256:dc42dfa52495c90dc5b99c19534d6d4fa9cd37fa439356fcbd73e770c35f2293"
        },
        {
            "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
            "size": 318908736,
            "digest": "sha256:841483099b542d6aeafc6ffacd59617954c409d3ebc558b7a95f43e05b1701a1"
        }
    ]
}`

// This is the same as the config digest.
const kindestImageDigest = `sha256:d6b20550c77b11385dd30115ba29dbf9a9bfc98c2f28ff7d162a6ad7c9686251`

const kindestConfig = `{
    "architecture": "amd64",
    "config": {
        "Hostname": "5e7483a6cf0e",
        "Domainname": "",
        "User": "",
        "AttachStdin": false,
        "AttachStdout": false,
        "AttachStderr": false,
        "Tty": false,
        "OpenStdin": false,
        "StdinOnce": false,
        "Env": [
            "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
            "container=docker",
            "HTTP_PROXY=",
            "HTTPS_PROXY=",
            "NO_PROXY="
        ],
        "Cmd": null,
        "Image": "docker.io/kindest/base:v20250521-31a79fd4",
        "Volumes": null,
        "WorkingDir": "/",
        "Entrypoint": [
            "/usr/local/bin/entrypoint",
            "/sbin/init"
        ],
        "OnBuild": null,
        "Labels": {},
        "StopSignal": "SIGRTMIN+3"
    },
    "container": "5e7483a6cf0e7958e796eee6912d1f6247394a0b914c822e47c7596b54aeac0a",
    "container_config": {
        "Hostname": "5e7483a6cf0e",
        "Domainname": "",
        "User": "",
        "AttachStdin": false,
        "AttachStdout": false,
        "AttachStderr": false,
        "Tty": false,
        "OpenStdin": false,
        "StdinOnce": false,
        "Env": [
            "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
            "container=docker"
        ],
        "Cmd": [
            "infinity"
        ],
        "Image": "docker.io/kindest/base:v20250521-31a79fd4",
        "Volumes": null,
        "WorkingDir": "/",
        "Entrypoint": [
            "sleep"
        ],
        "OnBuild": null,
        "Labels": {},
        "StopSignal": "SIGRTMIN+3"
    },
    "created": "2025-05-21T01:04:14.093628812Z",
    "docker_version": "20.10.21",
    "history": [
        {
            "created": "2025-05-21T00:57:52.347930888Z",
            "created_by": "COPY / / # buildkit",
            "comment": "buildkit.dockerfile.v0"
        },
        {
            "created": "2025-05-21T00:57:52.347930888Z",
            "created_by": "ENV container=docker",
            "comment": "buildkit.dockerfile.v0",
            "empty_layer": true
        },
        {
            "created": "2025-05-21T00:57:52.347930888Z",
            "created_by": "STOPSIGNAL SIGRTMIN+3",
            "comment": "buildkit.dockerfile.v0",
            "empty_layer": true
        },
        {
            "created": "2025-05-21T00:57:52.347930888Z",
            "created_by": "ENTRYPOINT [\"/usr/local/bin/entrypoint\" \"/sbin/init\"]",
            "comment": "buildkit.dockerfile.v0",
            "empty_layer": true
        },
        {
            "created": "2025-05-21T01:04:14.093628812Z",
            "created_by": "infinity"
        }
    ],
    "os": "linux",
    "rootfs": {
        "type": "layers",
        "diff_ids": [
            "sha256:f13bb3f5a0b612fa8b3ee54536cff9a57cc76a084b6a399d7762283e52393778",
            "sha256:2f6a2492037574ad66c92893c0d266ff1a74dbab0eeed7771c511a09f909b577"
        ]
    }
}`

func ImageSchemaV2(hooks ...func(*imageapi.Image)) *imageapi.Image {
	img := &imageapi.Image{
		ObjectMeta: metav1.ObjectMeta{
			Name:         kindestImageDigest,
			GenerateName: "kindest",
		},
		DockerImageReference:         "kindest/node:v1.33.1",
		DockerImageManifestMediaType: "application/vnd.docker.distribution.manifest.v2+json",
		DockerImageManifest:          kindestManifest,
		DockerImageConfig:            kindestConfig,
		DockerImageLayers: []imageapi.ImageLayer{
			{
				Name:      "sha256:d9d352c11bbd3880007953ed6eec1cbace76898828f3434984a0ca60672fdf5a",
				LayerSize: 29715337,
				MediaType: "application/vnd.oci.image.layer.v1.tar+gzip",
			},
		},
		DockerImageMetadata: imageapi.DockerImage{
			ID:            kindestImageDigest,
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
	for _, hook := range hooks {
		hook(img)
	}
	return img
}
