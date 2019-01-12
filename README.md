---

A simple tool to help manage credentials. Inspired by vault.

---

## Example

```bash
# Encrypt the secrets and save each in a file
> mkdir -p secret_store/postgres
> secmuxer -cmd encrypt -password p1ssw1rd > secret_store/postgres/host
pg
^d
> secmuxer -cmd encrypt -password p1ssw1rd > secret_store/postgres/username
readwrite
^d
> secmuxer -cmd encrypt -password p1ssw1rd > secret_store/postgres/password
s.cr.t
^d
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

## How did this come

- Where do I put my credentials?

That problem concerns me everytime I build a personal tool. 
I want everything to be version controlled, even the script of deployment.
And I also want my repo public.

I know vault for a long time. It runs as a service and needs an server.
But I just needs a simple cli tool. So I borrowed it's idea and made this.

- How secure is the encrypted data?
The code uses pbkdf2 (with SHA1 run 4096 times) to derive the encryption key from the password.
The actual encryption is using 256-bit AES-GCM with random 96-bit nonces copied from cryptopasta, see:
https://github.com/gtank/cryptopasta

