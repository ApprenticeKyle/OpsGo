# JWT密钥目录

请将以下文件放置在此目录：

- `public.pem` - RSA公钥文件
- `private_pkcs8.pem` - RSA私钥文件（PKCS8格式）

## 密钥格式说明

### 公钥格式（public.pem）
```pem
-----BEGIN PUBLIC KEY-----
...
-----END PUBLIC KEY-----
```

或

```pem
-----BEGIN RSA PUBLIC KEY-----
...
-----END RSA PUBLIC KEY-----
```

### 私钥格式（private_pkcs8.pem）
```pem
-----BEGIN PRIVATE KEY-----
...
-----END PRIVATE KEY-----
```

或

```pem
-----BEGIN RSA PRIVATE KEY-----
...
-----END RSA PRIVATE KEY-----
```

## 生成密钥对（如果需要）

```bash
# 生成私钥（PKCS8格式）
openssl genpkey -algorithm RSA -out private_pkcs8.pem -pkcs8 -pkeyopt rsa_keygen_bits:2048

# 从私钥提取公钥
openssl rsa -pubout -in private_pkcs8.pem -out public.pem
```
