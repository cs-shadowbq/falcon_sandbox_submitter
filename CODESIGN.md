# Code Signing

Objective:  Sign a binary with a self-signed certificate.

## Observer that we have no identities:

```shell
$> security find-identity -v -p codesigning
     0 valid identities found
$> security find-identity -p codesigning

Policy: Code Signing
  Matching identities
     0 identities found

  Valid identities only
     0 valid identities found
```

## Create a self-signed certificate:

Open Keychain Access.  
Choose Keychain Access > Certificate Assistant > Create Certificate ...  

    Enter a name 'cs-shadowbq'
    Set 'Certificate Type' to 'Code Signing'
    Set 'Let me override defaults' to 'Yes'
    Set 'Serial Number' to '1'
    Set 'Validity Period' to '3650' (10 years)
    Set 'Email' to 'me.at@example.com'

    Click 'Create'

## Observe that we have no valid identities.  We need to create a valid identity.

```shell
$> security find-identity -p codesigning

Policy: Code Signing
  Matching identities
  1) 9AC8280226A6B0D89C88874DCC9A15DA1B7312FD "cs-shadowbq" (CSSMERR_TP_NOT_TRUSTED)
     1 identities found

  Valid identities only
     0 valid identities found
$> security find-identity -v -p codesigning
     0 valid identities found
```

## Change the trust settings for the certificate:

Open the Keychain Access application and select the 'login' keychain.
Find the 'cs-shadowbq' certificate and double-click it.
Expand the 'Trust' section and set 'Code Signing' to 'Always Trust'.

## Now we have a valid identity:

```shell
$> security find-identity -p codesigning

Policy: Code Signing
  Matching identities
  1) 9AC8280226A6B0D89C88874DCC9A15DA1B7312FD "cs-shadowbq"
     1 identities found

  Valid identities only
  1) 9AC8280226A6B0D89C88874DCC9A15DA1B7312FD "cs-shadowbq"
     1 valid identities found

$> security find-identity -v -p codesigning
  1) 9AC8280226A6B0D89C88874DCC9A15DA1B7312FD "cs-shadowbq"
     1 valid identities found
```

## Sign the binary:

```shell
codesign -fs cs-shadowbq build/falcon_sandbox_darwin_amd64
codesign -fs cs-shadowbq build/falcon_sandbox.exe
```

## Verify the signature:

```shell
$> codesign --display -vvv "build/falcon_sandbox_darwin_amd64"
Executable=/Users/smacgregor/sandbox/falcon_sandbox/build/falcon_sandbox_darwin_amd64
Identifier=falcon_sandbox_darwin_amd64
Format=Mach-O thin (x86_64)
CodeDirectory v=20400 size=435444 flags=0x0(none) hashes=13602+2 location=embedded
Hash type=sha256 size=32
CandidateCDHash sha256=3b3b5eba531383370ba23b05766ef770e87db3a5
CandidateCDHashFull sha256=3b3b5eba531383370ba23b05766ef770e87db3a59510218f5e8cc241a4eceb6f
Hash choices=sha256
CMSDigest=3b3b5eba531383370ba23b05766ef770e87db3a59510218f5e8cc241a4eceb6f
CMSDigestType=2
CDHash=3b3b5eba531383370ba23b05766ef770e87db3a5
Signature size=1869
Authority=cs-shadowbq
Signed Time=Feb 17, 2024 at 8:56:18 PM
Info.plist=not bound
TeamIdentifier=not set
Sealed Resources=none
Internal requirements count=1 size=104
```

## Extract the public key from the certificate:

```shell
codesign --display --extract-certificates "build/falcon_sandbox_darwin_amd64"
mv codesign0 cs-shadowbq.pub.der
```

A pub x509 cert file (der formatted) will be create `codesign0` that contains the public key.

```shell
/usr/bin/openssl x509 -in cs-shadowbq.pub.der -inform der -text
Certificate:
    Data:
        Version: 3 (0x2)
        Serial Number: 1 (0x1)
    Signature Algorithm: sha256WithRSAEncryption
        Issuer: CN=cs-shadowbq, O=CrowdStrike, C=US/emailAddress=scott.macgregor@crowdstrike.com
        Validity
            Not Before: Feb 18 01:53:21 2024 GMT
            Not After : Jun 27 01:53:21 2027 GMT
        Subject: CN=cs-shadowbq, O=CrowdStrike, C=US/emailAddress=scott.macgregor@crowdstrike.com
        Subject Public Key Info:
            Public Key Algorithm: rsaEncryption
                RSA Public-Key: (2048 bit)
                Modulus:
                    00:e0:b6:e7:e1:3e:af:e0:5d:91:3d:04:6d:d8:59:
                    85:39:6a:80:c7:db:46:46:6c:fd:a5:70:54:07:93:
                    ac:d2:a2:6d:38:0a:70:81:e7:4d:4a:ee:7e:6e:79:
                    5c:7a:9b:ee:95:2e:ec:b4:37:08:d6:db:36:d9:90:
                    a1:6c:65:4d:ff:64:fc:de:be:fc:a2:18:40:4d:39:
                    32:23:c7:01:7c:f2:4b:49:ba:b2:80:26:8b:16:71:
                    26:d8:9c:f3:46:7f:a6:e4:c7:45:b5:a4:99:82:e4:
                    5c:8f:64:74:45:d6:a3:b5:77:0b:45:77:d5:84:c8:
                    98:68:bc:a1:b7:0d:55:2d:68:fb:ac:d5:4d:4f:d6:
                    77:29:2a:f1:04:b2:80:5b:5e:74:04:00:08:60:a1:
                    37:31:67:5f:a0:12:7a:51:04:b4:e2:e0:df:f5:59:
                    ff:45:c6:10:0d:66:f8:d7:59:79:2b:fb:3c:df:f1:
                    b9:f6:8f:00:3d:bc:31:58:a1:06:64:f2:9e:9e:cc:
                    a7:4d:c7:d9:cc:5e:ca:53:2e:f5:bc:6b:21:a4:5b:
                    87:a7:34:46:9f:55:8b:40:dc:60:63:88:db:8f:ca:
                    27:21:37:b4:e7:b4:ec:02:33:a1:4d:55:2c:22:78:
                    b8:70:61:5e:36:b4:85:7f:14:e8:00:7e:b4:56:93:
                    8a:9f
                Exponent: 65537 (0x10001)
        X509v3 extensions:
            X509v3 Key Usage: critical
                Digital Signature
            X509v3 Extended Key Usage: critical
                Code Signing
            X509v3 Subject Key Identifier:
                CA:7C:77:1D:0B:92:A3:A8:39:7A:82:19:23:26:B2:49:89:DC:7A:96
    Signature Algorithm: sha256WithRSAEncryption
         3d:cd:24:0b:12:2d:4d:d7:14:bb:71:b6:97:9f:12:0d:f4:7a:
         69:e7:aa:68:e4:4b:47:01:b4:0c:93:dc:d9:ff:83:bc:a9:75:
         cb:cc:60:0c:67:d6:89:c1:b6:c4:61:18:ae:e1:3a:5e:82:ac:
         8e:d5:1f:85:f7:24:c3:f3:4d:cc:d3:b1:43:47:48:4b:ee:73:
         36:f1:b1:71:60:7a:32:db:7e:1d:3a:46:1d:2a:4e:0b:78:eb:
         dd:af:8d:25:1f:ef:3f:3a:6a:2d:e6:b9:0f:04:cf:4e:48:68:
         01:5f:0a:59:08:b3:18:8c:1d:62:92:de:e2:e0:eb:e6:43:2a:
         e4:58:a2:5c:3a:7c:14:92:da:d2:44:12:3a:a9:89:c0:d9:71:
         ad:28:45:f2:50:6c:7b:0b:0b:34:8e:da:39:45:7e:f2:49:38:
         f2:2d:16:cc:42:4d:ad:fb:59:fc:e3:3d:13:fc:64:bc:b8:af:
         bf:7c:5c:00:4c:32:61:f7:02:dd:99:7a:7e:79:1f:89:ba:aa:
         e6:41:f9:74:9e:d5:3b:df:35:18:c9:27:99:c0:eb:84:ad:37:
         ff:82:7c:dd:be:28:42:8a:71:2d:06:a5:19:25:5d:69:ec:63:
         27:af:6c:73:d3:33:1e:5c:85:7c:b9:15:22:3e:e8:3d:a4:77:
         2e:4e:11:35
-----BEGIN CERTIFICATE-----
MIIDljCCAn6gAwIBAgIBATANBgkqhkiG9w0BAQsFADBpMRQwEgYDVQQDDAtjcy1z
aGFkb3dicTEUMBIGA1UECgwLQ3Jvd2RTdHJpa2UxCzAJBgNVBAYTAlVTMS4wLAYJ
KoZIhvcNAQkBFh9zY290dC5tYWNncmVnb3JAY3Jvd2RzdHJpa2UuY29tMB4XDTI0
MDIxODAxNTMyMVoXDTI3MDYyNzAxNTMyMVowaTEUMBIGA1UEAwwLY3Mtc2hhZG93
YnExFDASBgNVBAoMC0Nyb3dkU3RyaWtlMQswCQYDVQQGEwJVUzEuMCwGCSqGSIb3
DQEJARYfc2NvdHQubWFjZ3JlZ29yQGNyb3dkc3RyaWtlLmNvbTCCASIwDQYJKoZI
hvcNAQEBBQADggEPADCCAQoCggEBAOC25+E+r+BdkT0EbdhZhTlqgMfbRkZs/aVw
VAeTrNKibTgKcIHnTUrufm55XHqb7pUu7LQ3CNbbNtmQoWxlTf9k/N6+/KIYQE05
MiPHAXzyS0m6soAmixZxJtic80Z/puTHRbWkmYLkXI9kdEXWo7V3C0V31YTImGi8
obcNVS1o+6zVTU/Wdykq8QSygFtedAQACGChNzFnX6ASelEEtOLg3/VZ/0XGEA1m
+NdZeSv7PN/xufaPAD28MVihBmTynp7Mp03H2cxeylMu9bxrIaRbh6c0Rp9Vi0Dc
YGOI24/KJyE3tOe07AIzoU1VLCJ4uHBhXja0hX8U6AB+tFaTip8CAwEAAaNJMEcw
DgYDVR0PAQH/BAQDAgeAMBYGA1UdJQEB/wQMMAoGCCsGAQUFBwMDMB0GA1UdDgQW
BBTKfHcdC5KjqDl6ghkjJrJJidx6ljANBgkqhkiG9w0BAQsFAAOCAQEAPc0kCxIt
TdcUu3G2l58SDfR6aeeqaORLRwG0DJPc2f+DvKl1y8xgDGfWicG2xGEYruE6XoKs
jtUfhfckw/NNzNOxQ0dIS+5zNvGxcWB6Mtt+HTpGHSpOC3jr3a+NJR/vPzpqLea5
DwTPTkhoAV8KWQizGIwdYpLe4uDr5kMq5FiiXDp8FJLa0kQSOqmJwNlxrShF8lBs
ewsLNI7aOUV+8kk48i0WzEJNrftZ/OM9E/xkvLivv3xcAEwyYfcC3Zl6fnkfibqq
5kH5dJ7VO981GMknmcDrhK03/4J83b4oQopxLQalGSVdaexjJ69sc9MzHlyFfLkV
Ij7oPaR3Lk4RNQ==
-----END CERTIFICATE-----
```