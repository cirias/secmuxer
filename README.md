---

A simple tool to help manage credentials. Inspired by consul and vault.

---

## Example

```bash
> tree ./secret_store/
./secret_store/
└── postgres
    ├── host
    ├── username
    └── password

> secmuxer -store ./secret_store -password p1ssw1rd <<EOF
{
  "postgres": {
    "host": "{{ secret "postgres/host" }}",
    "port": "5432",
    "username": "{{ secret "postgres/username" }}",
    "password": "{{ secret "postgres/password" }}",
    "database": "db"
  }
}
EOF

# Output
{
  "postgres": {
    "host": "pg",
    "port": "5432",
    "username": "readwrite",
    "password": "s.cr.t",
    "database": "db"
  }
}
```

It's up to you how to store your secrets. They are just encrypted files. You can even version control them.

## Bad thing

1. The built in `sh` func is very dangerous. Only use this tool when you know the content of input.
2. To leverage simplicity, it runs `openssl` to decrypt the file. And you cannot change the encryption algorithm.
