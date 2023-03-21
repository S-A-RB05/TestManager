package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"example.com/m/v2/receive"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gorilla/mux"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func GetRunningContainers(w http.ResponseWriter, r *http.Request) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Printf("%s %s\n", container.ID[:10], container.Image)
	}
}

func CreateNewContainer(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	reader, err := cli.ImagePull(ctx, "docker.io/library/alpine", types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, reader)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "alpine",
		Cmd:   []string{"echo", "hello world"},
	}, nil, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
}

func StartContainer(w http.ResponseWriter, r *http.Request) {

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	if err := cli.ContainerStart(ctx, "0be914e5052a28459884fc535507751c57337e7a09a413c08c32b578b984b000", types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
}

func ExecuteCmd(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: cmd")

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	respIdExecCreate, err := cli.ContainerExecCreate(context.Background(), "0be914e5052a28459884fc535507751c57337e7a09a413c08c32b578b984b000", types.ExecConfig{
		User:       "root",
		Privileged: true,
		Cmd: []string{
			"sh", "-c", "wine terminal.exe",
		},
	})
	if err != nil {
		fmt.Println(err)
	}

	response, err := cli.ContainerExecAttach(context.Background(), respIdExecCreate.ID, types.ExecStartCheck{})
	if err != nil {
		panic(err)
	}
	defer response.Close()

	data, _ := ioutil.ReadAll(response.Reader)
	fmt.Println(string(data))
}

var b64 = `Ly8rLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tKwovL3wgICAgICAgICAgICAgICAgIEVBMzEzMzcgLSBtdWx0aS1zdHJhdGVneSBhZHZhbmNlZCB0cmF
kaW5nIHJvYm90LiB8Ci8vfCAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIENvcHlyaWdodCAyMDE2LTIwMjIsIEVBMzEzMzcgTHRkIHwKLy98ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgaHR0cHM6Ly9na
XRodWIuY29tL0VBMzEzMzcgfAovLystLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0rCgovKgogKiAgVGhpcyBmaWxlIGlzIGZyZWUgc29mdHdhcmU6IHlvdSBjYW4gcmV
kaXN0cmlidXRlIGl0IGFuZC9vciBtb2RpZnkKICogIGl0IHVuZGVyIHRoZSB0ZXJtcyBvZiB0aGUgR05VIEdlbmVyYWwgUHVibGljIExpY2Vuc2UgYXMgcHVibGlzaGVkIGJ5CiAqICB0aGUgRnJlZSBTb2Z0d2FyZSBGb3VuZGF0aW9uLCBla
XRoZXIgdmVyc2lvbiAzIG9mIHRoZSBMaWNlbnNlLCBvcgogKiAgKGF0IHlvdXIgb3B0aW9uKSBhbnkgbGF0ZXIgdmVyc2lvbi4KCiAqICBUaGlzIHByb2dyYW0gaXMgZGlzdHJpYnV0ZWQgaW4gdGhlIGhvcGUgdGhhdCBpdCB3aWxsIGJlIHV
zZWZ1bCwKICogIGJ1dCBXSVRIT1VUIEFOWSBXQVJSQU5UWTsgd2l0aG91dCBldmVuIHRoZSBpbXBsaWVkIHdhcnJhbnR5IG9mCiAqICBNRVJDSEFOVEFCSUxJVFkgb3IgRklUTkVTUyBGT1IgQSBQQVJUSUNVTEFSIFBVUlBPU0UuICBTZWUgd
GhlCiAqICBHTlUgR2VuZXJhbCBQdWJsaWMgTGljZW5zZSBmb3IgbW9yZSBkZXRhaWxzLgoKICogIFlvdSBzaG91bGQgaGF2ZSByZWNlaXZlZCBhIGNvcHkgb2YgdGhlIEdOVSBHZW5lcmFsIFB1YmxpYyBMaWNlbnNlCiAqICBhbG9uZyB3aXR
oIHRoaXMgcHJvZ3JhbS4gIElmIG5vdCwgc2VlIDxodHRwOi8vd3d3LmdudS5vcmcvbGljZW5zZXMvPi4KICovCgovLyBJbmNsdWRlcy4KI2luY2x1ZGUgImluY2x1ZGUvaW5jbHVkZXMuaCIKCi8vIEVBIHByb3BlcnRpZXMuCiNpZmRlZiBfX
3Byb3BlcnR5X18KI3Byb3BlcnR5IGNvcHlyaWdodCBlYV9jb3B5CiNwcm9wZXJ0eSBkZXNjcmlwdGlvbiBlYV9uYW1lCiNwcm9wZXJ0eSBkZXNjcmlwdGlvbiBlYV9kZXNjCiNwcm9wZXJ0eSBpY29uICJyZXNvdXJjZXMvZmF2aWNvbi5pY28
iCiNwcm9wZXJ0eSBsaW5rIGVhX2xpbmsKI3Byb3BlcnR5IHZlcnNpb24gZWFfdmVyc2lvbgojZW5kaWYKCi8vIEVBIGluZGljYXRvciByZXNvdXJjZXMuCiNpZmRlZiBfX3Jlc291cmNlX18KI2lmZGVmIF9fTVFMNV9fCi8vIFRlc3RlciBwc
m9wZXJ0aWVzLgojcHJvcGVydHkgdGVzdGVyX2luZGljYXRvciAiOjoiICsgSU5ESV9FV09fT1NDX1BBVEggKyAiXFxFbGxpb3R0X1dhdmVfT3NjaWxsYXRvcjIiICsgTVFMX0VYVAojcHJvcGVydHkgdGVzdGVyX2luZGljYXRvciAiOjoiICs
gSU5ESV9TVkVCQl9QQVRIICsgIlxcU1ZFX0JvbGxpbmdlcl9CYW5kcyIgKyBNUUxfRVhUCiNwcm9wZXJ0eSB0ZXN0ZXJfaW5kaWNhdG9yICI6OiIgKyBJTkRJX1RNQV9DR19QQVRIICsgIlxcVE1BK0NHX21sYWRlbl9OUlAiICsgTVFMX0VYV
AojcHJvcGVydHkgdGVzdGVyX2luZGljYXRvciAiOjoiICsgSU5ESV9BVFJfTUFfVFJFTkRfUEFUSCArICJcXEFUUl9NQV9UcmVuZCIgKyBNUUxfRVhUCiNwcm9wZXJ0eSB0ZXN0ZXJfaW5kaWNhdG9yICI6OiIgKyBJTkRJX1RNQV9UUlVFX1B
BVEggKyAiXFxUTUFfVHJ1ZSIgKyBNUUxfRVhUCiNwcm9wZXJ0eSB0ZXN0ZXJfaW5kaWNhdG9yICI6OiIgKyBJTkRJX1NBV0FfUEFUSCArICJcXFNBV0EiICsgTVFMX0VYVAojcHJvcGVydHkgdGVzdGVyX2luZGljYXRvciAiOjoiICsgSU5ES
V9TVVBFUlRSRU5EX1BBVEggKyAiXFxTdXBlclRyZW5kIiArIE1RTF9FWFQKLy8gSW5kaWNhdG9yIHJlc291cmNlcy4KI3Jlc291cmNlIElORElfRVdPX09TQ19QQVRIICsgIlxcRWxsaW90dF9XYXZlX09zY2lsbGF0b3IyIiArIE1RTF9FWFQ
KI3Jlc291cmNlIElORElfU1ZFQkJfUEFUSCArICJcXFNWRV9Cb2xsaW5nZXJfQmFuZHMiICsgTVFMX0VYVAojcmVzb3VyY2UgSU5ESV9UTUFfQ0dfUEFUSCArICJcXFRNQStDR19tbGFkZW5fTlJQIiArIE1RTF9FWFQKI3Jlc291cmNlIElOR
ElfQVRSX01BX1RSRU5EX1BBVEggKyAiXFxBVFJfTUFfVHJlbmQiICsgTVFMX0VYVAojcmVzb3VyY2UgSU5ESV9UTUFfVFJVRV9QQVRIICsgIlxcVE1BX1RydWUiICsgTVFMX0VYVAojcmVzb3VyY2UgSU5ESV9TQVdBX1BBVEggKyAiXFxTQVd
BIiArIE1RTF9FWFQKI3Jlc291cmNlIElORElfU1VQRVJUUkVORF9QQVRIICsgIlxcU3VwZXJUcmVuZCIgKyBNUUxfRVhUCiNlbmRpZgojZW5kaWYKCi8vIEdsb2JhbCB2YXJpYWJsZXMuCkVBICplYTsKCi8vIElucHV0IHZhcmlhYmxlcwoKa
W5wdXQgaW50IHRlc3R2YXIgPSAxOwppbnB1dCBib29sIHRlc3Rib29sID0gZmFsc2U7CmlucHV0IGRvdWJsZSB0ZXN0ZG91Yj0wLjE7CgovKiBFQSBldmVudCBoYW5kbGVyIGZ1bmN0aW9ucyAqLwoKLyoqCiAqIEluaXRpYWxpemF0aW9uIGZ
1bmN0aW9uIG9mIHRoZSBleHBlcnQuCiAqLwppbnQgT25Jbml0KCkKewogIGJvb2wgX2luaXRpYXRlZCA9IHRydWU7CiAgUHJpbnRGb3JtYXQoIiVzIHYlcyAoJXMpIGluaXRpYWxpemluZy4uLiIsIGVhX25hbWUsIGVhX3ZlcnNpb24sIGVhX
2xpbmspOwogIF9pbml0aWF0ZWQgJj0gSW5pdEVBKCk7CiAgX2luaXRpYXRlZCAmPSBJbml0U3RyYXRlZ2llcygpOwogIGlmIChHZXRMYXN0RXJyb3IoKSA+IDApCiAgewogICAgZWEuR2V0TG9nZ2VyKCkuRXJyb3IoIkVycm9yIGR1cmluZyB
pbml0aWFsaXppbmchIiwgX19GVU5DVElPTl9MSU5FX18sIFRlcm1pbmFsOjpHZXRMYXN0RXJyb3JUZXh0KCkpOwogIH0KICBpZiAoRUFfRGlzcGxheURldGFpbHNPbkNoYXJ0KQogIHsKICAgIERpc3BsYXlTdGFydHVwSW5mbyh0cnVlKTsKI
CB9CiAgZWEuR2V0TG9nZ2VyKCkuRmx1c2goKTsKICBDaGFydDo6V2luZG93UmVkcmF3KCk7CiAgaWYgKCFfaW5pdGlhdGVkKQogIHsKICAgIGVhLlNldChTVFJVQ1RfRU5VTShFQVN0YXRlLCBFQV9TVEFURV9GTEFHX0VOQUJMRUQpLCBmYWx
zZSk7CiAgfQogIHJldHVybiAoX2luaXRpYXRlZCA/IElOSVRfU1VDQ0VFREVEIDogSU5JVF9GQUlMRUQpOwp9CgovKioKICogRGVpbml0aWFsaXphdGlvbiBmdW5jdGlvbiBvZiB0aGUgZXhwZXJ0LgogKi8Kdm9pZCBPbkRlaW5pdChjb25zd
CBpbnQgcmVhc29uKSB7IERlaW5pdFZhcnMoKTsgfQoKLyoqCiAqICJUaWNrIiBldmVudCBoYW5kbGVyIGZ1bmN0aW9uIChFQSBvbmx5KS4KICoKICogSW52b2tlZCB3aGVuIGEgbmV3IHRpY2sgZm9yIGEgc3ltYm9sIGlzIHJlY2VpdmVkLCB
0byB0aGUgY2hhcnQgb2Ygd2hpY2ggdGhlIEV4cGVydCBBZHZpc29yIGlzIGF0dGFjaGVkLgogKi8Kdm9pZCBPblRpY2soKQp7CiAgRUFQcm9jZXNzUmVzdWx0IF9yZXN1bHQgPSBlYS5Qcm9jZXNzVGljaygpOwogIGlmIChfcmVzdWx0LnN0Z
19wcm9jZXNzZWRfcGVyaW9kcyA+IDApCiAgewogICAgaWYgKEVBX0Rpc3BsYXlEZXRhaWxzT25DaGFydCAmJiAoVGVybWluYWw6OklzVmlzdWFsTW9kZSgpIHx8IFRlcm1pbmFsOjpJc1JlYWx0aW1lKCkpKQogICAgewogICAgICBzdHJpbmc
gX3RleHQgPSBTdHJpbmdGb3JtYXQoIiVzIHYlcyBieSAlcyAoJXMpXG4iLCBlYV9uYW1lLCBlYV92ZXJzaW9uLCBlYV9hdXRob3IsIGVhX2xpbmspOwogICAgICBfdGV4dCArPQogICAgICAgICAgU2VyaWFsaXplckNvbnZlcnRlcjo6RnJvb
U9iamVjdChlYSwgU0VSSUFMSVpFUl9GTEFHX0lOQ0xVREVfRFlOQU1JQykuUHJlY2lzaW9uKDIpLlRvU3RyaW5nPFNlcmlhbGl6ZXJKc29uPigpOwogICAgICBfdGV4dCArPSBlYS5HZXRMb2dnZXIoKS5Ub1N0cmluZygpOwogICAgICBDb21t
ZW50KF90ZXh0KTsKICAgIH0KICB9Cn0KCiNpZmRlZiBfX01RTDVfXwovKioKICogIlRyYWRlIiBldmVudCBoYW5kbGVyIGZ1bmN0aW9uIChNUUw1IG9ubHkpLgogKgogKiBJbnZva2VkIHdoZW4gYSB0cmFkZSBvcGVyYXRpb24gaXMgY29tcGx
ldGVkIG9uIGEgdHJhZGUgc2VydmVyLgogKi8Kdm9pZCBPblRyYWRlKCkge30KCi8qKgogKiAiT25UcmFkZVRyYW5zYWN0aW9uIiBldmVudCBoYW5kbGVyIGZ1bmN0aW9uIChNUUw1IG9ubHkpLgogKgogKiBJbnZva2VkIHdoZW4gcGVyZm9ybW
luZyBzb21lIGRlZmluaXRlIGFjdGlvbnMgb24gYSB0cmFkZSBhY2NvdW50LCBpdHMgc3RhdGUgY2hhbmdlcy4KICovCnZvaWQgT25UcmFkZVRyYW5zYWN0aW9uKGNvbnN0IE1xbFRyYWRlVHJhbnNhY3Rpb24gJnRyYW5zLCAvLyBUcmFkZSB0c
mFuc2FjdGlvbiBzdHJ1Y3R1cmUuCiAgICAgICAgICAgICAgICAgICAgICAgIGNvbnN0IE1xbFRyYWRlUmVxdWVzdCAmcmVxdWVzdCwgICAvLyBSZXF1ZXN0IHN0cnVjdHVyZS4KICAgICAgICAgICAgICAgICAgICAgICAgY29uc3QgTXFsVHJh
ZGVSZXN1bHQgJnJlc3VsdCAgICAgIC8vIFJlc3VsdCBzdHJ1Y3R1cmUuCikKewp9CgovKioKICogIlRpbWVyIiBldmVudCBoYW5kbGVyIGZ1bmN0aW9uIChNUUw1IG9ubHkpLgogKgogKiBJbnZva2VkIHBlcmlvZGljYWxseSBnZW5lcmF0ZWQ
gYnkgdGhlIEVBIHRoYXQgaGFzIGFjdGl2YXRlZCB0aGUgdGltZXIgYnkgdGhlIEV2ZW50U2V0VGltZXIgZnVuY3Rpb24uCiAqIFVzdWFsbHksIHRoaXMgZnVuY3Rpb24gaXMgY2FsbGVkIGJ5IE9uSW5pdC4KICovCnZvaWQgT25UaW1lcigpIH
t9CgovKioKICogIlRlc3RlckluaXQiIGV2ZW50IGhhbmRsZXIgZnVuY3Rpb24gKE1RTDUgb25seSkuCiAqCiAqIFRoZSBzdGFydCBvZiBvcHRpbWl6YXRpb24gaW4gdGhlIHN0cmF0ZWd5IHRlc3RlciBiZWZvcmUgdGhlIGZpcnN0IG9wdGlta
XphdGlvbiBwYXNzLgogKgogKiBJbnZva2VkIHdpdGggdGhlIHN0YXJ0IG9mIG9wdGltaXphdGlvbiBpbiB0aGUgc3RyYXRlZ3kgdGVzdGVyLgogKgogKiBAc2VlOiBodHRwczovL3d3dy5tcWw1LmNvbS9lbi9kb2NzL2Jhc2lzL2Z1bmN0aW9u
L2V2ZW50cwogKi8Kdm9pZCBUZXN0ZXJJbml0KCkge30KCi8qKgogKiAiT25UZXN0ZXIiIGV2ZW50IGhhbmRsZXIgZnVuY3Rpb24uCiAqCiAqIEludm9rZWQgYWZ0ZXIgYSBoaXN0b3J5IHRlc3Rpbmcgb2YgYW4gRXhwZXJ0IEFkdmlzb3Igb24
gdGhlIGNob3NlbiBpbnRlcnZhbCBpcyBvdmVyLgogKiBJdCBpcyBjYWxsZWQgcmlnaHQgYmVmb3JlIHRoZSBjYWxsIG9mIE9uRGVpbml0KCkuCiAqCiAqIFJldHVybnMgY2FsY3VsYXRlZCB2YWx1ZSB0aGF0IGlzIHVzZWQgYXMgdGhlIEN1c3
RvbSBtYXggY3JpdGVyaW9uCiAqIGluIHRoZSBnZW5ldGljIG9wdGltaXphdGlvbiBvZiBpbnB1dCBwYXJhbWV0ZXJzLgogKgogKiBAc2VlOiBodHRwczovL3d3dy5tcWw1LmNvbS9lbi9kb2NzL2Jhc2lzL2Z1bmN0aW9uL2V2ZW50cwogKi8KL
y8gZG91YmxlIE9uVGVzdGVyKCkgeyByZXR1cm4gMS4wOyB9CgovKioKICogIk9uVGVzdGVyUGFzcyIgZXZlbnQgaGFuZGxlciBmdW5jdGlvbiAoTVFMNSBvbmx5KS4KICoKICogSW52b2tlZCB3aGVuIGEgZnJhbWUgaXMgcmVjZWl2ZWQgZHVy
aW5nIEV4cGVydCBBZHZpc29yIG9wdGltaXphdGlvbiBpbiB0aGUgc3RyYXRlZ3kgdGVzdGVyLgogKgogKiBAc2VlOiBodHRwczovL3d3dy5tcWw1LmNvbS9lbi9kb2NzL2Jhc2lzL2Z1bmN0aW9uL2V2ZW50cwogKi8Kdm9pZCBPblRlc3RlclB
hc3MoKSB7fQoKLyoqCiAqICJPblRlc3RlckRlaW5pdCIgZXZlbnQgaGFuZGxlciBmdW5jdGlvbiAoTVFMNSBvbmx5KS4KICoKICogSW52b2tlZCBhZnRlciB0aGUgZW5kIG9mIEV4cGVydCBBZHZpc29yIG9wdGltaXphdGlvbiBpbiB0aGUgc3
RyYXRlZ3kgdGVzdGVyLgogKgogKiBAc2VlOiBodHRwczovL3d3dy5tcWw1LmNvbS9lbi9kb2NzL2Jhc2lzL2Z1bmN0aW9uL2V2ZW50cwogKi8Kdm9pZCBPblRlc3RlckRlaW5pdCgpIHt9CgovKioKICogIk9uQm9va0V2ZW50IiBldmVudCBoY
W5kbGVyIGZ1bmN0aW9uIChNUUw1IG9ubHkpLgogKgogKiBJbnZva2VkIG9uIERlcHRoIG9mIE1hcmtldCBjaGFuZ2VzLgogKiBUbyBwcmUtc3Vic2NyaWJlIHVzZSB0aGUgTWFya2V0Qm9va0FkZCgpIGZ1bmN0aW9uLgogKiBJbiBvcmRlciB0
byB1bnN1YnNjcmliZSBmb3IgYSBwYXJ0aWN1bGFyIHN5bWJvbCwgY2FsbCBNYXJrZXRCb29rUmVsZWFzZSgpLgogKi8Kdm9pZCBPbkJvb2tFdmVudChjb25zdCBzdHJpbmcgJnN5bWJvbCkge30KCi8qKgogKiAiT25Cb29rRXZlbnQiIGV2ZW5
0IGhhbmRsZXIgZnVuY3Rpb24gKE1RTDUgb25seSkuCiAqCiAqIEludm9rZWQgYnkgdGhlIGNsaWVudCB0ZXJtaW5hbCB3aGVuIGEgdXNlciBpcyB3b3JraW5nIHdpdGggYSBjaGFydC4KICovCnZvaWQgT25DaGFydEV2ZW50KGNvbnN0IGludC
BpZCwgICAgICAgICAvLyBFdmVudCBJRC4KICAgICAgICAgICAgICAgICAgY29uc3QgbG9uZyAmbHBhcmFtLCAgIC8vIFBhcmFtZXRlciBvZiB0eXBlIGxvbmcgZXZlbnQuCiAgICAgICAgICAgICAgICAgIGNvbnN0IGRvdWJsZSAmZHBhcmFtL
CAvLyBQYXJhbWV0ZXIgb2YgdHlwZSBkb3VibGUgZXZlbnQuCiAgICAgICAgICAgICAgICAgIGNvbnN0IHN0cmluZyAmc3BhcmFtICAvLyBQYXJhbWV0ZXIgb2YgdHlwZSBzdHJpbmcgZXZlbnRzLgopCnsKfQoKLy8gQHRvZG86IE9uVHJhZGVU
cmFuc2FjdGlvbiAoaHR0cHM6Ly93d3cubXFsNS5jb20vZW4vZG9jcy9iYXNpcy9mdW5jdGlvbi9ldmVudHMpLgojZW5kaWYgLy8gZW5kOiBfX01RTDVfXwoKLyogQ3VzdG9tIEVBIGZ1bmN0aW9ucyAqLwoKLyoqCiAqIERpc3BsYXkgaW5mbyB
vbiB0aGUgY2hhcnQuCiAqLwpib29sIERpc3BsYXlTdGFydHVwSW5mbyhib29sIF9zdGFydHVwID0gZmFsc2UsIHN0cmluZyBzZXAgPSAiXG4iKQp7CiAgc3RyaW5nIF9vdXRwdXQgPSAiIjsKICBSZXNldExhc3RFcnJvcigpOwogIGlmIChlYS
5HZXRTdGF0ZSgpLklzT3B0aW1pemF0aW9uTW9kZSgpIHx8IChlYS5HZXRTdGF0ZSgpLklzVGVzdGluZ01vZGUoKSAmJiAhZWEuR2V0U3RhdGUoKS5Jc1Zpc3VhbE1vZGUoKSkpCiAgewogICAgLy8gSWdub3JlIGNoYXJ0IHVwZGF0ZXMgd2hlb
iBvcHRpbWl6aW5nIG9yIHRlc3RpbmcgaW4gbm9uLXZpc3VhbCBtb2RlLgogICAgcmV0dXJuIGZhbHNlOwogIH0KICBfb3V0cHV0ICs9ICJBQ0NPVU5UOiAiICsgZWEuQWNjb3VudCgpLlRvU3RyaW5nKCkgKyBzZXA7CiAgX291dHB1dCArPSAi
RUE6ICIgKyBlYS5Ub1N0cmluZygpICsgc2VwOwogIF9vdXRwdXQgKz0gIlRFUk1JTkFMOiAiICsgZWEuR2V0VGVybWluYWwoKS5Ub1N0cmluZygpICsgc2VwOwojaWZkZWYgX19hZHZhbmNlZF9fCiAgLy8gUHJpbnQgZW5hYmxlZCBzdHJhdGV
naWVzIGluZm8uCiAgZm9yIChEaWN0U3RydWN0SXRlcmF0b3I8bG9uZywgUmVmPFN0cmF0ZWd5Pj4gX3NpdGVyID0gZWEuR2V0U3RyYXRlZ2llcygpLkJlZ2luKCk7IF9zaXRlci5Jc1ZhbGlkKCk7ICsrX3NpdGVyKQogIHsKICAgIFN0cmF0ZW
d5ICpfc3RyYXQgPSBfc2l0ZXIuVmFsdWUoKS5QdHIoKTsKICAgIHN0cmluZyBfc25hbWUgPQogICAgICAgIF9zdHJhdC5HZXROYW1lKCk7IC8vICsgIkAiICsgQ2hhcnRUZjo6VGZUb1N0cmluZyhfc3RyYXQuR2V0VGYoKS5HZXQ8RU5VTV9US
U1FRlJBTUVTPihDSEFSVF9QQVJBTV9URikpOwogICAgX291dHB1dCArPSBTdHJpbmdGb3JtYXQoIlN0cmF0ZWd5OiAlczogJXNcbiIsIF9zbmFtZSwKICAgICAgICAgICAgICAgICAgICAgICAgICAgIFNlcmlhbGl6ZXJDb252ZXJ0ZXI6OkZy
b21PYmplY3QoX3N0cmF0LCBTRVJJQUxJWkVSX0ZMQUdfSU5DTFVERV9EWU5BTUlDKQogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIC5Ub1N0cmluZzxTZXJpYWxpemVySnNvbj4oU0VSSUFMSVpFUl9KU09OX05PX1dISVRFU1BBQ0V
TKSk7CiAgfQojZW5kaWYKICBpZiAoX3N0YXJ0dXApCiAgewogICAgaWYgKGVhLkdldChTVFJVQ1RfRU5VTShFQVN0YXRlLCBFQV9TVEFURV9GTEFHX1RSQURFX0FMTE9XRUQpKSkKICAgIHsKICAgICAgaWYgKCFUZXJtaW5hbDo6SGFzRXJyb3
IoKSkKICAgICAgewogICAgICAgIF9vdXRwdXQgKz0gc2VwICsgIlRyYWRpbmcgaXMgYWxsb3dlZCwgd2FpdGluZyBmb3IgbmV3IGJhcnMuLi4iOwogICAgICB9CiAgICAgIGVsc2UKICAgICAgewogICAgICAgIF9vdXRwdXQgKz0gc2VwICsgI
lRyYWRpbmcgaXMgYWxsb3dlZCwgYnV0IHRoZXJlIGlzIHNvbWUgaXNzdWUuLi4iOwogICAgICAgIF9vdXRwdXQgKz0gc2VwICsgVGVybWluYWw6OkdldExhc3RFcnJvclRleHQoKTsKICAgICAgICBlYS5HZXRMb2dnZXIoKS5BZGRMYXN0RXJy
b3IoX19GVU5DVElPTl9MSU5FX18pOwogICAgICB9CiAgICB9CiAgICBlbHNlIGlmIChUZXJtaW5hbDo6SXNSZWFsdGltZSgpKQogICAgewogICAgICBfb3V0cHV0ICs9IHNlcCArIFN0cmluZ0Zvcm1hdCgKICAgICAgICAgICAgICAgICAgICA
gICAgICAgIkVycm9yICVkOiBUcmFkaW5nIGlzIG5vdCBhbGxvd2VkIGZvciB0aGlzIHN5bWJvbCwgcGxlYXNlIGVuYWJsZSBhdXRvbWF0ZWQgdHJhZGluZyBvciBjaGVjayAiCiAgICAgICAgICAgICAgICAgICAgICAgICAgICJ0aGUgc2V0dG
luZ3MhIiwKICAgICAgICAgICAgICAgICAgICAgICAgICAgX19MSU5FX18pOwogICAgfQogICAgZWxzZQogICAgewogICAgICBfb3V0cHV0ICs9IHNlcCArICJXYWl0aW5nIGZvciBuZXcgYmFycy4uLiI7CiAgICB9CiAgfQogIENvbW1lbnQoX
291dHB1dCk7CiAgcmV0dXJuICFUZXJtaW5hbDo6SGFzRXJyb3IoKTsKfQoKLyoqCiAqIEluaXQgRUEuCiAqLwpib29sIEluaXRFQSgpCnsKICBib29sIF9pbml0aWF0ZWQgPSBlYV9hdXRoOwogIEVBUGFyYW1zIGVhX3BhcmFtcyhfX0ZJTEVf
XywgVmVyYm9zZUxldmVsKTsKICAvLyBlYV9wYXJhbXMuU2V0Q2hhcnRJbmZvRnJlcShFQV9EaXNwbGF5RGV0YWlsc09uQ2hhcnQgPyAyIDogMCk7CiAgLy8gRUEgcGFyYW1zLgogIGVhX3BhcmFtcy5TZXREZXRhaWxzKGVhX25hbWUsIGVhX2R
lc2MsIGVhX3ZlcnNpb24sIFN0cmluZ0Zvcm1hdCgiJXMgKCVzKSIsIGVhX2F1dGhvciwgZWFfbGluaykpOwogIC8vIFJpc2sgcGFyYW1zLgogIGVhX3BhcmFtcy5TZXQoU1RSVUNUX0VOVU0oRUFQYXJhbXMsIEVBX1BBUkFNX1BST1BfUklTS1
9NQVJHSU5fTUFYKSwgRUFfUmlza19NYXJnaW5NYXgpOwogIGVhX3BhcmFtcy5TZXRGbGFnKEVBX1BBUkFNX0ZMQUdfTE9UU0laRV9BVVRPLCBFQV9Mb3RTaXplIDw9IDApOwogIC8vIEluaXQgaW5zdGFuY2UuCiAgZWEgPSBuZXcgRUEoZWFfc
GFyYW1zKTsKICBlYS5TZXQoVFJBREVfUEFSQU1fUklTS19NQVJHSU4sIEVBX1Jpc2tfTWFyZ2luTWF4KTsKICBpZiAoIWVhLkdldChTVFJVQ1RfRU5VTShFQVN0YXRlLCBFQV9TVEFURV9GTEFHX1RSQURFX0FMTE9XRUQpKSkKICB7CiAgICBl
YS5HZXRMb2dnZXIoKS5FcnJvcigKICAgICAgICAiVHJhZGluZyBpcyBub3QgYWxsb3dlZCBmb3IgdGhpcyBzeW1ib2wsIHBsZWFzZSBlbmFibGUgYXV0b21hdGVkIHRyYWRpbmcgb3IgY2hlY2sgdGhlIHNldHRpbmdzISIsCiAgICAgICAgX19
GVU5DVElPTl9MSU5FX18pOwogICAgX2luaXRpYXRlZCAmPSBmYWxzZTsKICB9CiNpZmRlZiBfX2FkdmFuY2VkX18KICBpZiAoX2luaXRpYXRlZCkKICB7CiAgICBFQVRhc2tzIF9lYV90YXNrcyhlYSk7CiAgICBfaW5pdGlhdGVkICY9IF9lYV90YXNrcy5BZGRUYXNrKEVBX1Rhc2sxX0lmLCBFQV9UYXNrMV9UaGVuKTsKICAgIF9pbml0aWF0ZWQgJj0gX2VhX3Rhc2tzLkFk
ZFRhc2soRUFfVGFzazJfSWYsIEVBX1Rhc2syX1RoZW4pOwogICAgX2luaXRpYXRlZCAmPSBfZWFfdGFza3MuQWRkVGFzayhFQV9UYXNrM19JZiwgRUFfVGFzazNfVGhlbik7CiAgfQojZW5kaWYKICByZXR1cm4gX2luaXRpYXRlZDsKfQoKLyoqCiAqIEluaXQgc3RyYXRlZ2llcy4KICovCmJvb2wgSW5pdFN0cmF0ZWdpZXMoKQp7CiAgYm9vbCBfcmVzID0gZWFfZXhpc3RzOwogIGludCBfbWFnaWNfc3RlcCA9IEZJTkFMX0VOVU1fVElNRUZSQU1FU19JTkRFWDsKICBsb25nIF9tYWdpY19ubyA9IEVBX01hZ2ljTnVtYmVyOwogIFJlc2V0TGFzdEVycm9yKCk7CiAgLy8gSW5pdGlhbGl6ZSBzdHJhdGVnaWVzIHBlciB0aW1lZnJhbWUuCiAgRUFTdHJhdGVneUFkZChTdHJhdGVneV9NMSwgMSA8PCBNMSk7CiAgRUFTdHJhdGVneUFkZChTdHJhdGVneV9NNSwgMSA8PCBNNSk7CiAgRUFTdHJhdGVneUFkZChTdHJhdGVneV9NMTUsIDEgPDwgTTE1KTsKICBFQVN0cmF0ZWd5QWRkKFN0cmF0ZWd5X00zMCwgMSA8PCBNMzApOwogIEVBU3RyYXRlZ3lBZGQoU3RyYXRlZ3lfSDEsIDEgPDwgSDEpOwogIEVBU3RyYXRlZ3lBZGQoU3RyYXRlZ3lfSDIsIDEgPDwgSDIpOwogIEVBU3RyYXRlZ3lBZGQoU3RyYXRlZ3lfSDMsIDEgPDwgSDMpOwogIEVBU3RyYXRlZ3lBZGQoU3RyYXRlZ3lfSDQsIDEgPDwgSDQpOwogIEVBU3RyYXRlZ3lBZGQoU3RyYXRlZ3lfSDYsIDEgPDwgSDYpOwogIEVBU3RyYXRlZ3lBZGQoU3RyYXRlZ3lfSDgsIDEgPDwgSDgpOwogIEVBU3RyYXRlZ3lBZGQoU3RyYXRlZ3lfSDEyLCAxIDw8IEgxMik7CiAgLy8gVXBkYXRlIGxvdCBzaXplLgogIGVhLlNldChTVFJBVF9QQVJBTV9MUywgRUFfTG90U2l6ZSk7CiAgLy8gT3ZlcnJpZGUgbWF4IHNwcmVhZCB2YWx1ZXMuCiAgZWEuU2V0KFNUUkFUX1BBUkFNX01BWF9TUFJFQUQsIEVBX01heFNwcmVhZCk7CiAgZWEuU2V0KFRSQURFX1BBUkFNX01BWF9TUFJFQUQsIEVBX01heFNwcmVhZCk7CiNpZmRlZiBfX2FkdmFuY2VkX18KICBlYS5TZXQoU1RSQVRfUEFSQU1fU09GTSwgRUFfU2lnbmFsT3BlbkZpbHRlck1ldGhvZCk7CiAgZWEuU2V0KFNUUkFUX1BBUkFNX1NDRk0sIEVBX1NpZ25hbENsb3NlRmlsdGVyTWV0aG9kKTsKICBlYS5TZXQoU1RSQVRfUEFSQU1fU09GVCwgRUFfU2lnbmFsT3BlbkZpbHRlclRpbWUpOwogIGVhLlNldChTVFJBVF9QQVJBTV9URk0sIEVBX1RpY2tGaWx0ZXJNZXRob2QpOwogIGVhLlNldChTVFJVQ1RfRU5VTShFQVBhcmFtcywgRUFfUEFSQU1fUFJPUF9TSUdOQUxfRklMVEVSKSwgRUFfU2lnbmFsT3BlblN0cmF0ZWd5RmlsdGVyKTsKI2lmZGVmIF9fcmlkZXJfXwogIC8vIERpc2FibGVzIHN0cmF0ZWd5IGRlZmluZWQgb3JkZXIgY2xvc3VyZXMgZm9yIFJpZGVyLgogIGVhLlNldChTVFJBVF9QQVJBTV9PQ0wsIDApOwogIGVhLlNldChTVFJBVF9QQVJBTV9PQ1AsIDApOwogIGVhLlNldChTVFJBVF9QQVJBTV9PQ1QsIDApOwogIC8vIEluaXQgcHJpY2Ugc3RvcCBtZXRob2RzIGZvciBhbGwgdGltZWZyYW1lcy4KICBFQVN0cmF0ZWd5QWRkU3RvcHMoTlVMTCwgRUFfU3RvcHNfU3RyYXQsIEVBX1N0b3BzX1RmKTsKI2Vsc2UKICBlYS5TZXQoU1RSQVRfUEFSQU1fT0NMLCBFQV9PcmRlckNsb3NlTG9zcyk7CiAgZWEuU2V0KFNUUkFUX1BBUkFNX09DUCwgRUFfT3JkZXJDbG9zZVByb2ZpdCk7CiAgZWEuU2V0KFNUUkFUX1BBUkFNX09DVCwgRUFfT3JkZXJDbG9zZVRpbWUpOwogIC8vIEluaXQgcHJpY2Ugc3RvcCBtZXRob2RzIGZvciBlYWNoIHRpbWVmcmFtZS4KICBfcmVzICY9IEVBU3RyYXRlZ3lBZGRTdG9wcyhlYS5HZXRTdHJhdGVneVZpYVByb3A8aW50PihTVFJBVF9QQVJBTV9URiwgUEVSSU9EX00xKSwgRUFfU3RvcHNfTTEsIFBFUklPRF9NMSk7CiAgX3JlcyAmPSBFQVN0cmF0ZWd5QWRkU3RvcHMoZWEuR2V0U3RyYXRlZ3lWaWFQcm9wPGludD4oU1RSQVRfUEFSQU1fVEYsIFBFUklPRF9NNSksIEVBX1N0b3BzX001LCBQRVJJT0RfTTUpOwogIF9yZXMgJj0gRUFTdHJhdGVneUFkZFN0b3BzKGVhLkdldFN0cmF0ZWd5VmlhUHJvcDxpbnQ+KFNUUkFUX1BBUkFNX1RGLCBQRVJJT0RfTTE1KSwgRUFfU3RvcHNfTTE1LCBQRVJJT0RfTTE1KTsKICBfcmVzICY9IEVBU3RyYXRlZ3lBZGRTdG9wcyhlYS5HZXRTdHJhdGVneVZpYVByb3A8aW50PihTVFJBVF9QQVJBTV9URiwgUEVSSU9EX00zMCksIEVBX1N0b3BzX00zMCwgUEVSSU9EX00zMCk7CiAgX3JlcyAmPSBFQVN0cmF0ZWd5QWRkU3RvcHMoZWEuR2V0U3RyYXRlZ3lWaWFQcm9wPGludD4oU1RSQVRfUEFSQU1fVEYsIFBFUklPRF9IMSksIEVBX1N0b3BzX0gxLCBQRVJJT0RfSDEpOwogIF9yZXMgJj0gRUFTdHJhdGVneUFkZFN0b3BzKGVhLkdldFN0cmF0ZWd5VmlhUHJvcDxpbnQ+KFNUUkFUX1BBUkFNX1RGLCBQRVJJT0RfSDIpLCBFQV9TdG9wc19IMiwgUEVSSU9EX0gyKTsKICBfcmVzICY9IEVBU3RyYXRlZ3lBZGRTdG9wcyhlYS5HZXRTdHJhdGVneVZpYVByb3A8aW50PihTVFJBVF9QQVJBTV9URiwgUEVSSU9EX0gzKSwgRUFfU3RvcHNfSDMsIFBFUklPRF9IMyk7CiAgX3JlcyAmPSBFQVN0cmF0ZWd5QWRkU3RvcHMoZWEuR2V0U3RyYXRlZ3lWaWFQcm9wPGludD4oU1RSQVRfUEFSQU1fVEYsIFBFUklPRF9INCksIEVBX1N0b3BzX0g0LCBQRVJJT0RfSDQpOwogIF9yZXMgJj0gRUFTdHJhdGVneUFkZFN0b3BzKGVhLkdldFN0cmF0ZWd5VmlhUHJvcDxpbnQ+KFNUUkFUX1BBUkFNX1RGLCBQRVJJT0RfSDYpLCBFQV9TdG9wc19INiwgUEVSSU9EX0g2KTsKICBfcmVzICY9IEVBU3RyYXRlZ3lBZGRTdG9wcyhlYS5HZXRTdHJhdGVneVZpYVByb3A8aW50PihTVFJBVF9QQVJBTV9URiwgUEVSSU9EX0g4KSwgRUFfU3RvcHNfSDgsIFBFUklPRF9IOCk7CiAgX3JlcyAmPSBFQVN0cmF0ZWd5QWRkU3RvcHMoZWEuR2V0U3RyYXRlZ3lWaWFQcm9wPGludD4oU1RSQVRfUEFSQU1fVEYsIFBFUklPRF9IMTIpLCBFQV9TdG9wc19IMTIsIFBFUklPRF9IMTIpOwojZW5kaWYgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAvLyBfX3JpZGVyX18KI2VuZGlmICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgLy8gX19hZHZhbmNlZF9fCiAgX3JlcyAmPSBHZXRMYXN0RXJyb3IoKSA9PSAwIHx8IEdldExhc3RFcnJvcigpID09IDUwNTM7IC8vIEBmaXhtZTogZXJyb3IgNTA1Mz8KICBSZXNldExhc3RFcnJvcigpOwogIHJldHVybiBfcmVzICYmIGVhX2NvbmZpZ3VyZWQ7Cn0KCi8qKgogKiBBZGRzIHN0cmF0ZWd5IHRvIHRoZSBnaXZlbiB0aW1lZnJhbWUuCiAqLwpib29sIEVBU3RyYXRlZ3lBZGQoRU5VTV9TVFJBVEVHWSBfc3RnLCBpbnQgX3RmcykKewogIGJvb2wgX3Jlc3VsdCA9IHRydWU7CiAgdW5zaWduZWQgaW50IF9tYWdpY19ubyA9IEVBX01hZ2ljTnVtYmVyICsgX3N0ZyAqIEZJTkFMX0VOVU1fVElNRUZSQU1FU19JTkRFWDsKICBzd2l0Y2ggKF9zdGcpCiAgewogIGNhc2UgU1RSQVRfQUM6CiAgICByZXR1cm4gZWEuU3RyYXRlZ3lBZGQ8U3RnX0FDPihfdGZzLCBfbWFnaWNfbm8sIF9zdGcpOwogIGNhc2UgU1RSQVRfQUQ6CiAgICByZXR1cm4gZWEuU3RyYXRlZ3lBZGQ8U3RnX0FEPihfdGZzLCBfbWFnaWNfbm8sIF9zdGcpOwogIGNhc2UgU1RSQVRfQURYOgogICAgcmV0dXJuIGVhLlN0cmF0ZWd5QWRkPFN0Z19BRFg+KF90ZnMsIF9tYWdpY19ubywgX3N0Zyk7CiAgY2FzZSBTVFJBVF9BTUE6CiAgICByZXR1cm4gZWEuU3RyYXRlZ3lBZGQ8U3RnX0FNQT4oX3RmcywgX21hZ2ljX25vLCBfc3RnKTsKICBjYXNlIFNUUkFUX0FTSToKICAgIHJldHVybiBlYS5TdHJhdGVneUFkZDxTdGdfQVNJPihfdGZzLCBfbWFnaWNfbm8sIF9zdGcpOwogIGNhc2UgU1RSQVRfQVRSOgogICAgcmV0dXJuIGVhLlN0cmF0ZWd5QWRkPFN0Z19BVFI+KF90ZnMsIF9tYWdpY19ubywgX3N0Zyk7CiAgY2FzZSBTVFJBVF9BTExJR0FUT1I6CiAgICByZXR1cm4gZWEuU3RyYXRlZ3lBZGQ8U3RnX0FsbGlnYXRvcj4oX3RmcywgX21hZ2ljX25vLCBfc3RnKTsKICBjYXNlIFNUUkFUX0FXRVNPTUU6CiAgICByZXR1cm4gZWEuU3RyYXRlZ3lBZGQ8U3RnX0F3ZXNvbWU+KF90ZnMsIF9tYWdpY19ubywgX3N0Zyk7CiNpZmRlZiBfX01RTDVfXwogICAgLy8gY2FzZSBTVFJBVF9BVFJfTUFfVFJFTkQ6CiAgICAvLyByZXR1cm4gZWEuU3RyYXRlZ3lBZGQ8U3RnX0FUUl9NQV9UcmVuZD4oX3RmcywgX21hZ2ljX25vLCBfc3RnKTsKI2VuZGlmCiAgY2FzZSBTVFJBVF9CV01GSToKICAgIHJldHVybiBlYS5TdHJhdGVneUFkZDxTdGdfQldNRkk+KF90ZnMsIF9tYWdpY19ubywgX3N0Zyk7CiAgY2FzZSBTVFJBVF9CQU5EUzoKICAgIHJldHVybiBlYS5TdHJhdGVneUFkZDxTdGdfQmFuZHM+KF90ZnMsIF9tYWdpY19ubywgX3N0Zyk7CiAgY2FzZSBTVFJBVF9CRUFSU19QT1dFUjoKICAgIHJldHVybiBlYS5TdHJhdGVneUFkZDxTdGdfQmVhcnNQb3dlcj4oX3RmcywgX21hZ2ljX25vLCBfc3RnKTsKICBjYXNlIFNUUkFUX0JVTExTX1BPV0VSOgogICAgcmV0dXJuIGVhLlN0cmF0ZWd5QWRkPFN0Z19CdWxsc1Bvd2VyPihfdGZzLCBfbWFnaWNfbm8sIF9zdGcpOwogIGNhc2UgU1RSQVRfQ0NJOgogICAgcmV0dXJuIGVhLlN0cmF0ZWd5QWRkPFN0Z19DQ0k+KF90ZnMsIF9tYWdpY19ubywgX3N0Zyk7CiAgY2FzZSBTVFJBVF9DSEFJS0lOOgogICAgcmV0dXJuIGVhLlN0cmF0ZWd5QWRkPFN0Z19DaGFpa2luPihfdGZzLCBfbWFnaWNfbm8sIF9zdGcpOwogIGNhc2UgU1RSQVRfREVNQToKICAgIHJldHVybiBlYS5TdHJhdGVneUFkZDxTdGdfREVNQT4oX3RmcywgX21hZ2ljX25vLCBfc3RnKTsKICBjYXNlIFNUUkFUX0RFTUFSS0VSOgogICAgcmV0dXJuIGVhLlN0cmF0ZWd5QWRkPFN0Z19EZU1hcmtlcj4oX3RmcywgX21hZ2ljX25vLCBfc3RnKTsKICBjYXNlIFNUUkFUX0VOVkVMT1BFUzoKICAgIHJldHVybiBlYS5TdHJhdGVneUFkZDxTdGdfRW52ZWxvcGVzPihfdGZzLCBfbWFnaWNfbm8sIF9zdGcpOwogIGNhc2UgU1RSQVRfRVdPOgogICAgcmV0dXJuIGVhLlN0cmF0ZWd5QWRkPFN0Z19FbGxpb3R0V2F2ZT4oX3RmcywgX21hZ2ljX25vLCBfc3RnKTsKICBjYXNlIFNUUkFUX0ZPUkNFOgogICAgcmV0dXJuIGVhLlN0cmF0ZWd5QWRkPFN0Z19Gb3JjZT4oX3RmcywgX21hZ2ljX25vLCBfc3RnKTsKICBjYXNlIFNUUkFUX0ZSQUNUQUxTOgogICAgcmV0dXJuIGVhLlN0cmF0ZWd5QWRkPFN0Z19GcmFjdGFscz4oX3RmcywgX21hZ2ljX25vLCBfc3RnKTsKICBjYXNlIFNUUkFUX0dBVE9SOgogICAgcmV0dXJuIGVhLlN0cmF0ZWd5QWRkPFN0Z19HYXRvcj4oX3RmcywgX21hZ2ljX25vLCBfc3RnKTsKICBjYXNlIFNUUkFUX0hFSUtFTl9BU0hJOgogICAgcmV0dXJuIGVhLlN0cmF0ZWd5QWRkPFN0Z19IZWlrZW5Bc2hpPihfdGZzLCBfbWFnaWNfbm8sIF9zdGcpOwogIGNhc2UgU1RSQVRfSUNISU1PS1U6CiAgICByZXR1cm4gZWEuU3RyYXRlZ3lBZGQ8U3RnX0ljaGltb2t1PihfdGZzLCBfbWFnaWNfbm8sIF9zdGcpOwogIGNhc2UgU1RSQVRfTUE6CiAgICByZXR1cm4gZWEuU3RyYXRlZ3lBZGQ8U3RnX01BPihfdGZzLCBfbWFnaWNfbm8sIF9zdGcpOwogIGNhc2UgU1RSQVRfTUFDRDoKICAgIHJldHVybiBlYS5TdHJhdGVneUFkZDxTdGdfTUFDRD4oX3RmcywgX21hZ2ljX25vLCBfc3RnKTsKICBjYXNlIFNUUkFUX01GSToKICAgIHJldHVybiBlYS5TdHJhdGVneUFkZDxTdGdfTUZJPihfdGZzLCBfbWFnaWNfbm8sIF9zdGcpOwogIGNhc2UgU1RSQVRfTU9NRU5UVU06CiAgICByZXR1cm4gZWEuU3RyYXRlZ3lBZGQ8U3RnX01vbWVudHVtPihfdGZzLCBfbWFnaWNfbm8sIF9zdGcpOwogIGNhc2UgU1RSQVRfT0JWOgogICAgcmV0dXJuIGVhLlN0cmF0ZWd5QWRkPFN0Z19PQlY+KF90ZnMsIF9tYWdpY19ubywgX3N0Zyk7CiAgY2FzZSBTVFJBVF9PU01BOgogICAgcmV0dXJuIGVhLlN0cmF0ZWd5QWRkPFN0Z19Pc01BPihfdGZzLCBfbWFnaWNfbm8sIF9zdGcpOwogIGNhc2UgU1RSQVRfUEFUVEVSTjoKICAgIHJldHVybiBlYS5TdHJhdGVneUFkZDxTdGdfUGF0dGVybj4oX3RmcywgX21hZ2ljX25vLCBfc3RnKTsKICBjYXNlIFNUUkFUX1BJTkJBUjoKICAgIHJldHVybiBlYS5TdHJhdGVneUFkZDxTdGdfUGluYmFyPihfdGZzLCBfbWFnaWNfbm8sIF9zdGcpOwogIGNhc2UgU1RSQVRfUElWT1Q6CiAgICByZXR1cm4gZWEuU3RyYXRlZ3lBZGQ8U3RnX1Bpdm90PihfdGZzLCBfbWFnaWNfbm8sIF9zdGcpOwogIGNhc2UgU1RSQVRfUlNJOgogICAgcmV0dXJuIGVhLlN0cmF0ZWd5QWRkPFN0Z19SU0k+KF90ZnMsIF9tYWdpY19ubywgX3N0Zyk7CiAgY2FzZSBTVFJBVF9SVkk6CiAgICByZXR1cm4gZWEuU3RyYXRlZ3lBZGQ8U3RnX1JWST4oX3RmcywgX21hZ2ljX25vLCBfc3RnKTsKICBjYXNlIFNUUkFUX1NBUjoKICAgIHJldHVybiBlYS5TdHJhdGVneUFkZDxTdGdfU0FSPihfdGZzLCBfbWFnaWNfbm8sIF9zdGcpOwogIC8vIGNhc2UgU1RSQVRfU0FXQToKICAvLyByZXR1cm4gZWEuU3RyYXRlZ3lBZGQ8U3RnX1NBV0E+KF90ZnMsIF9tYWdpY19ubywgX3N0Zyk7CiAgY2FzZSBTVFJBVF9TVERERVY6CiAgICByZXR1cm4gZWEuU3RyYXRlZ3lBZGQ8U3RnX1N0ZERldj4oX3RmcywgX21hZ2ljX25vLCBfc3RnKTsKICBjYXNlIFNUUkFUX1NUT0NIQVNUSUM6CiAgICByZXR1cm4gZWEuU3RyYXRlZ3lBZGQ8U3RnX1N0b2NoYXN0aWM+KF90ZnMsIF9tYWdpY19ubywgX3N0Zyk7CiNpZmRlZiBfX01RTDVfXwogICAgLy8gY2FzZSBTVFJBVF9TVVBFUlRSRU5EOgogICAgLy8gcmV0dXJuIGVhLlN0cmF0ZWd5QWRkPFN0Z19TdXBlclRyZW5kPihfdGZzLCBfbWFnaWNfbm8sIF9zdGcpOwojZW5kaWYKICBjYXNlIFNUUkFUX1NWRV9CQjoKICAgIHJldHVybiBlYS5TdHJhdGVneUFkZDxTdGdfU1ZFX0JvbGxpbmdlcl9CYW5kcz4oX3RmcywgX21hZ2ljX25vLCBfc3RnKTsKICBjYXNlIFNUUkFUX1RNQVRfU1ZFQkI6CiAgICByZXR1cm4gZWEuU3RyYXRlZ3lBZGQ8U3RnX1RNQVRfU1ZFQkI+KF90ZnMsIF9tYWdpY19ubywgX3N0Zyk7CiAgLy8gY2FzZSBTVFJBVF9UTUFfQ0c6CiAgLy8gX3Jlc3VsdCAmPSBlYS5TdHJhdGVneUFkZDxTdGdfVE1BX0NHPihfdGZzLCBfbWFnaWNfbm8sIF9zdGcpOwogIGNhc2UgU1RSQVRfVE1BX1RSVUU6CiAgICBfcmVzdWx0ICY9IGVhLlN0cmF0ZWd5QWRkPFN0Z19UTUFfVHJ1ZT4oX3RmcywgX21hZ2ljX25vLCBfc3RnKTsKICAgIC8vIF9yZXN1bHQgJj0gZWEuR2V0U3RyYXRlZ3koX21hZ2ljX25vKS5TZXQoKTsgQHRvZG8KICAgIGJyZWFrOwogIGNhc2UgU1RSQVRfV1BSOgogICAgcmV0dXJuIGVhLlN0cmF0ZWd5QWRkPFN0Z19XUFI+KF90ZnMsIF9tYWdpY19ubywgX3N0Zyk7CiAgY2FzZSBTVFJBVF9aSUdaQUc6CiAgICByZXR1cm4gZWEuU3RyYXRlZ3lBZGQ8U3RnX1ppZ1phZz4oX3RmcywgX21hZ2ljX25vLCBfc3RnKTsKICBjYXNlIFNUUkFUX05PTkU6CiAgICBicmVhazsKICBkZWZhdWx0OgogICAgU2V0VXNlckVycm9yKEVSUl9JTlZBTElEX1BBUkFNRVRFUik7CiAgICBfcmVzdWx0ICY9IGZhbHNlOwogICAgYnJlYWs7CiAgfQogIHJldHVybiBfcmVzdWx0Owp9CgovKioKICogQWRkcyBzdHJhdGVneSBzdG9wcy4KICovCmJvb2wgRUFTdHJhdGVneUFkZFN0b3BzKFN0cmF0ZWd5ICpfc3RyYXQgPSBOVUxMLCBFTlVNX1NUUkFURUdZIF9lbnVtX3N0Z19zdG9wcyA9IFNUUkFUX05PTkUsIEVOVU1fVElNRUZSQU1FUyBfdGYgPSAwKQp7CiAgYm9vbCBfcmVzdWx0ID0gdHJ1ZTsKICBpZiAoX2VudW1fc3RnX3N0b3BzID09IFNUUkFUX05PTkUpCiAgewogICAgcmV0dXJuIF9yZXN1bHQ7CiAgfQogIFN0cmF0ZWd5ICpfc3RyYXRfc3RvcHMgPSBlYS5HZXRTdHJhdGVneVZpYVByb3AyPGludCwgaW50PihTVFJBVF9QQVJBTV9UWVBFLCBfZW51bV9zdGdfc3RvcHMsIFNUUkFUX1BBUkFNX1RGLCBfdGYpOwogIGlmICghX3N0cmF0X3N0b3BzKQogIHsKICAgIF9yZXN1bHQgJj0gRUFTdHJhdGVneUFkZChfZW51bV9zdGdfc3RvcHMsIDEgPDwgQ2hhcnRUZjo6VGZUb0luZGV4KF90ZikpOwogICAgX3N0cmF0X3N0b3BzID0gZWEuR2V0U3RyYXRlZ3lWaWFQcm9wMjxpbnQsIGludD4oU1RSQVRfUEFSQU1fVFlQRSwgX2VudW1fc3RnX3N0b3BzLCBTVFJBVF9QQVJBTV9URiwgX3RmKTsKICAgIGlmIChfc3RyYXRfc3RvcHMpCiAgICB7CiAgICAgIF9zdHJhdF9zdG9wcy5FbmFibGVkKGZhbHNlKTsKICAgIH0KICB9CiAgaWYgKF9zdHJhdF9zdG9wcykKICB7CiAgICBpZiAoX3N0cmF0ICE9IE5VTEwgJiYgX3RmID4gMCkKICAgIHsKICAgICAgX3N0cmF0LlNldFN0b3BzKF9zdHJhdF9zdG9wcywgX3N0cmF0X3N0b3BzKTsKICAgIH0KICAgIGVsc2UKICAgIHsKICAgICAgZm9yIChEaWN0U3RydWN0SXRlcmF0b3I8bG9uZywgUmVmPFN0cmF0ZWd5Pj4gaXRlciA9IGVhLkdldFN0cmF0ZWdpZXMoKS5CZWdpbigpOyBpdGVyLklzVmFsaWQoKTsgKytpdGVyKQogICAgICB7CiAgICAgICAgU3RyYXRlZ3kgKl9zdHJhdF9yZWYgPSBpdGVyLlZhbHVlKCkuUHRyKCk7CiAgICAgICAgaWYgKF9zdHJhdF9yZWYuSXNFbmFibGVkKCkpCiAgICAgICAgewogICAgICAgICAgX3N0cmF0X3JlZi5TZXRTdG9wcyhfc3RyYXRfc3RvcHMsIF9zdHJhdF9zdG9wcyk7CiAgICAgICAgfQogICAgICB9CiAgICB9CiAgfQogIHJldHVybiBfcmVzdWx0Owp9CgovKioKICogRGVpbml0aWFsaXplIGdsb2JhbCBjbGFzcyB2YXJpYWJsZXMuCiAqLwp2b2lkIERlaW5pdFZhcnMoKSB7IE9iamVjdDo6RGVsZXRlKGVhKTsgfQo=`

func DecodeBase64(w http.ResponseWriter, r *http.Request) {
	dec, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		panic(err)
	}

	f, err := os.Create("myfilename.ex5")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err := f.Write(dec); err != nil {
		panic(err)
	}
	if err := f.Sync(); err != nil {
		panic(err)
	}

	// go to begginng of file
	f.Seek(0, 0)

	// output file contents
	io.Copy(os.Stdout, f)
}

func GenerateConfig() {
	// Create a new file to write the configuration settings to
	file, err := os.Create("config.ini")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// Write the configuration settings to the file
	_, err = file.WriteString("[Common]\n")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = file.WriteString("Login=123456\n")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = file.WriteString("Password=password\n")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print a success message
	fmt.Println("config.ini file generated successfully")
}

func handleRequests() {
	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)
	// replace http.HandleFunc with myRouter.HandleFunc
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/running", GetRunningContainers)
	myRouter.HandleFunc("/create", CreateNewContainer)
	myRouter.HandleFunc("/start", StartContainer)
	myRouter.HandleFunc("/cmd", ExecuteCmd)
	myRouter.HandleFunc("/decode", DecodeBase64)

	// finally, instead of passing in nil, we want
	// to pass in our newly created router as the second
	// argument
	log.Fatal(http.ListenAndServe(":8081", myRouter))
}

func main() {
	handleRequests()
	receive.Strat()
}
