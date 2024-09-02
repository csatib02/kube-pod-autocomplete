# Security

## Verifying container images

To verify KPA container-images you can use the following public key:

```key
-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEk3vpOBc1zsCdQ+vU56tETv64F5RW
ISazzb8sOyUqrkKV/JRe7Xb0OnaqGY7KopsIIxbrX+CbyCdQDtN73qf5EA==
-----END PUBLIC KEY-----
```

Save the key to a file:

```bash
echo '-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEk3vpOBc1zsCdQ+vU56tETv64F5RW
ISazzb8sOyUqrkKV/JRe7Xb0OnaqGY7KopsIIxbrX+CbyCdQDtN73qf5EA==
-----END PUBLIC KEY-----' > cosign.pub
```

Verify an image:

```bash
cosign verify --key cosign.pub <image url>
```
