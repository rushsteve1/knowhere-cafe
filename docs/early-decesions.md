# Why Go?

Why not?
Go is extremely well suited for writing networking applications.
It's fast *enough*, simple to use, has top class tooling, and please just
[read the list of std pakages](https://pkg.go.dev/std).

Plus I'm very familiar with it and I've used it a lot.

# Why Postgres?

I agonized about this for longer than I should.
I **really** wanted to use
[modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite)
which is a library I really like.

But the deciding factor was that Postgres has more features,
is even *more* robust, and I'm already connecting to SMTP, IMAP,
and potentially several other services.
Doubling down on PG early will save me a lot of hassle later.

# Why Passwordless?

Oh how I wish I could go back in time and remake the web
[a la Hypnospace](https://store.steampowered.com/app/844590/Hypnospace_Outlaw/).

Passwords were never a good idea.
They're a poor imitation of proper cryptography practices.
So instead I chose to rely on newer browser features such as Passkeys,
or falling back to Magic Links sent via email.

# Why not JMAP?

I did look into this, but there's just not a lot of support for it anywhere.
I'd be locked into a small handful of email providers/software.

# Dependencies

See [`go.mod`](../go.mod)

## Gorm


