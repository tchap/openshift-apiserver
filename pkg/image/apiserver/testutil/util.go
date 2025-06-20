package testutil

import (
	"fmt"
)

// InternalRegistryURL is an url of internal container image registry for testing purposes.
const InternalRegistryURL = "172.30.12.34:5000"

// MakeDockerImageReference makes a container image reference string referencing testing internal docker
// registry.
func MakeDockerImageReference(ns, isName, imageID string) string {
	return fmt.Sprintf("%s/%s/%s@%s", InternalRegistryURL, ns, isName, imageID)
}

// BaseImageWith1LayerDigest is the digest associated with BaseImageWith1Layer.
//
// This is actually docksal/empty.
const BaseImageWith1LayerDigest = `sha256:f853843b26903da94dd1cdf9e39ff7e2ba7a754388341895d557dbe913f5a915`

// BaseImageWith1Layer contains a single layer.
const BaseImageWith1Layer = `{
   "schemaVersion": 2,
   "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
   "config": {
      "mediaType": "application/vnd.docker.container.image.v1+json",
      "size": 1512,
      "digest": "sha256:6c6084ed97e5851b5d216b20ed1852301278584c3c6aff915272b231593f6f98"
   },
   "layers": [
      {
         "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
         "size": 1970140,
         "digest": "sha256:550fe1bea624a5c62551cf09f3aa10886eed133794844af1dfb775118309387e"
      }
   ]
}`

// BaseImageWith1LayerConfig is the config associated with BaseImageWith1Layer.
const BaseImageWith1LayerConfig = `{
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
      "/bin/sh"
    ],
    "ArgsEscaped": true,
    "Image": "sha256:01b85c6717c3b1f5379864199c541cecabb81be758b4fec6ef0b66cbfb6e11a5",
    "Volumes": null,
    "WorkingDir": "",
    "Entrypoint": null,
    "OnBuild": null,
    "Labels": null
  },
  "container": "f8a4df32c288f30c6d641c3945c88b64490e1e029be516209955023786cf1727",
  "container_config": {
    "Hostname": "f8a4df32c288",
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
      "CMD [\"/bin/sh\"]"
    ],
    "ArgsEscaped": true,
    "Image": "sha256:01b85c6717c3b1f5379864199c541cecabb81be758b4fec6ef0b66cbfb6e11a5",
    "Volumes": null,
    "WorkingDir": "",
    "Entrypoint": null,
    "OnBuild": null,
    "Labels": {}
  },
  "created": "2018-01-09T21:13:01.402230769Z",
  "docker_version": "17.06.2-ce",
  "history": [
    {
      "created": "2018-01-09T21:13:01.165340448Z",
      "created_by": "/bin/sh -c #(nop) ADD file:df48d6d6df42a01380557aebd4ca02807fc08a76a1d1b36d957e59a41c69db0b in / "
    },
    {
      "created": "2018-01-09T21:13:01.402230769Z",
      "created_by": "/bin/sh -c #(nop)  CMD [\"/bin/sh\"]",
      "empty_layer": true
    }
  ],
  "os": "linux",
  "rootfs": {
    "type": "layers",
    "diff_ids": [
      "sha256:d39d92664027be502c35cf1bf464c726d15b8ead0e3084be6e252a161730bc82"
    ]
  }
}`

// BaseImageWith2LayersDigest is the digest associated with BaseImageWith2Layers.
const BaseImageWith2LayersDigest = "sha256:6c64418d228763c6f01f81c209efa169f9930a691724d8b10faf025fc6ece53f"

// BaseImageWith2Layers contains 2 layers while the first one is shared with BaseImageWith1Layer.
const BaseImageWith2Layers = `{
  "schemaVersion": 2,
  "mediaType": "application/vnd.oci.image.manifest.v1+json",
  "config": {
    "mediaType": "application/vnd.oci.image.config.v1+json",
    "digest": "sha256:245780beb82f97496ad0dee2d70c6fd2bfa56adedb229a3e25c9e801dd262b5e",
    "size": 2803
  },
  "layers": [
    {
      "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
      "digest": "sha256:155ad54a8b2812a0ec559ff82c0c6f0f0dddb337a226b11879f09e15f67b69fc",
      "size": 48476100
    },
    {
      "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
      "digest": "sha256:8031108f3cda87bb32f090262d0109c8a0db99168050967becefad502e9a681b",
      "size": 24058530
    }
  ]
}`

// based on baseImageWith1Layer, it adds a new data layer of 126 B
const ChildImageWith2LayersDigest = "sha256:a9f073fbf2c9835711acd09081d87f5b7129ac6269e0df834240000f48abecd4"

const ChildImageWith2Layers = `{
   "schemaVersion": 1,
   "name": "miminar/childImageWith2Layers",
   "tag": "latest",
   "architecture": "amd64",
   "fsLayers": [
      {
         "blobSum": "sha256:766b6e9134dc2819fae9c5e67d39e14272948bc8967df9a119418cca84cab089"
      },
      {
         "blobSum": "sha256:a3ed95caeb02ffe68cdd9fd84406680ae93d633cb16422d00e8a7c22955b46d4"
      },
      {
         "blobSum": "sha256:2d099e04ef6c850542d8ab916df2e9417cc799d39b78f64440e51402f1261a36"
      },
      {
         "blobSum": "sha256:a3ed95caeb02ffe68cdd9fd84406680ae93d633cb16422d00e8a7c22955b46d4"
      }
   ],
   "history": [
      {
          "v1Compatibility": "{\"architecture\":\"amd64\",\"author\":\"miminar@redhat.com\",\"config\":{\"Hostname\":\"d7b63ae1152b\",\"Domainname\":\"\",\"User\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":[\"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\"],\"Cmd\":null,\"Image\":\"sha256:27bc5bf237c48c2b41b0636a3876960a9adb6c2ac9ff95ac879d56b1046ba5a1\",\"Volumes\":null,\"WorkingDir\":\"\",\"Entrypoint\":null,\"OnBuild\":[],\"Labels\":{}},\"container\":\"c2d2505e43f4fd479aa21d356270d0791633e838284d7010cba1f61992907c69\",\"container_config\":{\"Hostname\":\"d7b63ae1152b\",\"Domainname\":\"\",\"User\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":[\"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\"],\"Cmd\":[\"/bin/sh\",\"-c\",\"#(nop) COPY file:859e4175fd5743f276905245e351272b425232cfd3b30a3fc6bff351da308996 in /data3\"],\"Image\":\"sha256:27bc5bf237c48c2b41b0636a3876960a9adb6c2ac9ff95ac879d56b1046ba5a1\",\"Volumes\":null,\"WorkingDir\":\"\",\"Entrypoint\":null,\"OnBuild\":[],\"Labels\":{}},\"created\":\"2016-02-15T07:33:17.59074814Z\",\"docker_version\":\"1.10.0\",\"id\":\"e6a8e2793d6cad7d503aa5a3b55fd2c19b3b190d480a175b21d5f7b50c86d27b\",\"os\":\"linux\",\"parent\":\"84dc393745ff2631760c4bdbf1168af188fcd4606c1400c6900487fdc75a9ed5\",\"size\":126}"
      },
      {
         "v1Compatibility": "{\"id\":\"84dc393745ff2631760c4bdbf1168af188fcd4606c1400c6900487fdc75a9ed5\",\"parent\":\"1620fdccc2424391c3422467cec611bc32767d5bfae5bd8a2fb53c795e2a3e86\",\"created\":\"2016-02-15T07:33:17.454934648Z\",\"container_config\":{\"Cmd\":[\"/bin/sh -c #(nop) MAINTAINER miminar@redhat.com\"]},\"throwaway\":true}"
      },
      {
         "v1Compatibility": "{\"id\":\"1620fdccc2424391c3422467cec611bc32767d5bfae5bd8a2fb53c795e2a3e86\",\"parent\":\"3690474eb5b4b26fdfbd89c6e159e8cc376ca76ef48032a30fa6aafd56337880\",\"created\":\"2016-02-15T07:30:37.655693399Z\",\"container_config\":{\"Cmd\":[\"/bin/sh -c #(nop) COPY file:90583fd8c765e40f7f2070c55da446e138b019b0712dee898d8193b66b05d48d in /data1\"]},\"size\":128}"
      },
      {
         "v1Compatibility": "{\"id\":\"3690474eb5b4b26fdfbd89c6e159e8cc376ca76ef48032a30fa6aafd56337880\",\"created\":\"2016-02-15T07:30:37.531741167Z\",\"container_config\":{\"Cmd\":[\"/bin/sh -c #(nop) MAINTAINER miminar@redhat.com\"]},\"throwaway\":true}"
      }
   ]
}`

// based on baseImageWith2Layers, it adds a new data layer of 70 B
const ChildImageWith3LayersDigest = "sha256:2282a6d553353756fa43ba8672807d3fe81f8fdef54b0f6a360d64aaef2f243a"

const ChildImageWith3Layers = `{
   "schemaVersion": 1,
   "name": "miminar/childImageWith3Layers",
   "tag": "latest",
   "architecture": "amd64",
   "fsLayers": [
      {
         "blobSum": "sha256:77ef66f4abb43c5e17bcacdfe744f6959365f6244b66a6565470083fbdd15178"
      },
      {
         "blobSum": "sha256:a3ed95caeb02ffe68cdd9fd84406680ae93d633cb16422d00e8a7c22955b46d4"
      },
      {
         "blobSum": "sha256:e7900a2e6943680b384950859a0616089757cae4d8c6e98db9cfec6c41fe2834"
      },
      {
         "blobSum": "sha256:2d099e04ef6c850542d8ab916df2e9417cc799d39b78f64440e51402f1261a36"
      },
      {
         "blobSum": "sha256:a3ed95caeb02ffe68cdd9fd84406680ae93d633cb16422d00e8a7c22955b46d4"
      }
   ],
   "history": [
      {
         "v1Compatibility": "{\"architecture\":\"amd64\",\"author\":\"miminar@redhat.com\",\"config\":{\"Hostname\":\"686b99d75c4a\",\"Domainname\":\"\",\"User\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":[\"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\"],\"Cmd\":null,\"Image\":\"sha256:8b0241d44c66c1bcf48c66d0465ee6bf6ac2117e9936a9ec2337122e08d109ef\",\"Volumes\":null,\"WorkingDir\":\"\",\"Entrypoint\":null,\"OnBuild\":[],\"Labels\":{}},\"container\":\"61c9522f27b7052081b61b72d70dd71ce7050566812f050158e03954b493e446\",\"container_config\":{\"Hostname\":\"686b99d75c4a\",\"Domainname\":\"\",\"User\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":[\"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\"],\"Cmd\":[\"/bin/sh\",\"-c\",\"#(nop) COPY file:7781db9ed3a36b0607009b073a99802a9ad834bbb5e3bcbcf83a7d27146a1a5b in /data4\"],\"Image\":\"sha256:8b0241d44c66c1bcf48c66d0465ee6bf6ac2117e9936a9ec2337122e08d109ef\",\"Volumes\":null,\"WorkingDir\":\"\",\"Entrypoint\":null,\"OnBuild\":[],\"Labels\":{}},\"created\":\"2016-02-15T07:36:13.703778299Z\",\"docker_version\":\"1.10.0\",\"id\":\"8e7b1ec73ed1d21747991c2101d1db51e97c4f62931bbaa575aeba11286d6748\",\"os\":\"linux\",\"parent\":\"fbe31426cd0e8c5545ddc5c8318499682d52ff96118e36e49616ac3aee32c47c\",\"size\":70}"
      },
      {
         "v1Compatibility": "{\"id\":\"fbe31426cd0e8c5545ddc5c8318499682d52ff96118e36e49616ac3aee32c47c\",\"parent\":\"9b1154060650718a3850e625464addb217c1064f18dd693cf635dfcabdc9de50\",\"created\":\"2016-02-15T07:36:13.585345649Z\",\"container_config\":{\"Cmd\":[\"/bin/sh -c #(nop) MAINTAINER miminar@redhat.com\"]},\"throwaway\":true}"
      },
      {
         "v1Compatibility": "{\"id\":\"9b1154060650718a3850e625464addb217c1064f18dd693cf635dfcabdc9de50\",\"parent\":\"1620fdccc2424391c3422467cec611bc32767d5bfae5bd8a2fb53c795e2a3e86\",\"created\":\"2016-02-15T07:31:50.390272025Z\",\"container_config\":{\"Cmd\":[\"/bin/sh -c #(nop) COPY file:23d2e6ff1c67ff4caee900c71d58df6e37bfb9defe46085018c4ba29c3d2de5a in /data2\"]},\"size\":112}"
      },
      {
         "v1Compatibility": "{\"id\":\"1620fdccc2424391c3422467cec611bc32767d5bfae5bd8a2fb53c795e2a3e86\",\"parent\":\"3690474eb5b4b26fdfbd89c6e159e8cc376ca76ef48032a30fa6aafd56337880\",\"created\":\"2016-02-15T07:30:37.655693399Z\",\"container_config\":{\"Cmd\":[\"/bin/sh -c #(nop) COPY file:90583fd8c765e40f7f2070c55da446e138b019b0712dee898d8193b66b05d48d in /data1\"]},\"size\":128}"
      },
      {
         "v1Compatibility": "{\"id\":\"3690474eb5b4b26fdfbd89c6e159e8cc376ca76ef48032a30fa6aafd56337880\",\"created\":\"2016-02-15T07:30:37.531741167Z\",\"container_config\":{\"Cmd\":[\"/bin/sh -c #(nop) MAINTAINER miminar@redhat.com\"]},\"throwaway\":true}"
      }
   ]
}`

// another base image with unique data layer of 554 B
const MiscImageDigest = "sha256:2643199e5ed5047eeed22da854748ed88b3a63ba0497601ba75852f7b92d4640"

const MiscImage = `{
   "schemaVersion": 1,
   "name": "miminar/misc",
   "tag": "latest",
   "architecture": "amd64",
   "fsLayers": [
      {
         "blobSum": "sha256:a3ed95caeb02ffe68cdd9fd84406680ae93d633cb16422d00e8a7c22955b46d4"
      },
      {
         "blobSum": "sha256:eeee0535bf3cec7a24bff2c6e97481afa3d37e2cdeff277c57cb5cbdb2fa9e92"
      }
   ],
   "history": [
      {
         "v1Compatibility": "{\"id\":\"964092b7f3e54185d3f425880be0b022bfc9a706701390e0ceab527c84dea3e3\",\"parent\":\"9e77fef7a1c9f989988c06620dabc4020c607885b959a2cbd7c2283c91da3e33\",\"created\":\"2016-01-15T18:06:41.282540103Z\",\"container\":\"4e937d31f242d087cce0ec5b9fdbceaf1a13b40704e9147962cc80947e4ab86b\",\"container_config\":{\"Hostname\":\"aded96b43f48\",\"Domainname\":\"\",\"User\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":null,\"Cmd\":[\"/bin/sh\",\"-c\",\"#(nop) CMD [\\\"sh\\\"]\"],\"Image\":\"9e77fef7a1c9f989988c06620dabc4020c607885b959a2cbd7c2283c91da3e33\",\"Volumes\":null,\"WorkingDir\":\"\",\"Entrypoint\":null,\"OnBuild\":null,\"Labels\":null},\"docker_version\":\"1.8.3\",\"config\":{\"Hostname\":\"aded96b43f48\",\"Domainname\":\"\",\"User\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":null,\"Cmd\":[\"sh\"],\"Image\":\"9e77fef7a1c9f989988c06620dabc4020c607885b959a2cbd7c2283c91da3e33\",\"Volumes\":null,\"WorkingDir\":\"\",\"Entrypoint\":null,\"OnBuild\":null,\"Labels\":null},\"architecture\":\"amd64\",\"os\":\"linux\"}"
      },
      {
         "v1Compatibility": "{\"id\":\"9e77fef7a1c9f989988c06620dabc4020c607885b959a2cbd7c2283c91da3e33\",\"created\":\"2016-01-15T18:06:40.707908287Z\",\"container\":\"aded96b43f48d94eb80642c210b89f119ab2a233c1c7c7055104fb052937f12c\",\"container_config\":{\"Hostname\":\"aded96b43f48\",\"Domainname\":\"\",\"User\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":null,\"Cmd\":[\"/bin/sh\",\"-c\",\"#(nop) ADD file:a62b361be92f978752150570261ddc6fc21b025e3a28418820a1f39b7db7498c in /\"],\"Image\":\"\",\"Volumes\":null,\"WorkingDir\":\"\",\"Entrypoint\":null,\"OnBuild\":null,\"Labels\":null},\"docker_version\":\"1.8.3\",\"config\":{\"Hostname\":\"aded96b43f48\",\"Domainname\":\"\",\"User\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":null,\"Cmd\":null,\"Image\":\"\",\"Volumes\":null,\"WorkingDir\":\"\",\"Entrypoint\":null,\"OnBuild\":null,\"Labels\":null},\"architecture\":\"amd64\",\"os\":\"linux\",\"Size\":554}"
      }
   ]
}`

const ManifestList = `{
  "manifests": [
    {
      "digest": "sha256:96a76fa48db5fca24271fe1565d88a4453e759b365dbaaeeb5a4e41049293e77",
      "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
      "platform": {
        "architecture": "amd64",
        "os": "linux"
      },
      "size": 429
    },
    {
      "digest": "sha256:6c3d8fec1c50ff78997e13a8352b030d4b290f656081c974373753fd5a3496f1",
      "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
      "platform": {
        "architecture": "arm64",
        "os": "linux"
      },
      "size": 429
    },
    {
      "digest": "sha256:520a368f78807947b96ea773cc62b14e380f4af08bbfd8ed18f0ebc70dedef68",
      "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
      "platform": {
        "architecture": "ppc64le",
        "os": "linux"
      },
      "size": 429
    },
    {
      "digest": "sha256:50b0c55990fe1b48c4b026fb6b49b4377e36c52b291d434977793dc0c8998ba4",
      "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
      "platform": {
        "architecture": "s390x",
        "os": "linux"
      },
      "size": 429
    }
  ],
  "mediaType": "application/vnd.docker.distribution.manifest.list.v2+json",
  "schemaVersion": 2
}`
