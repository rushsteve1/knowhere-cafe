# Why Go?

Why not?
Go is extremely well suited for writing networking applications.
It's fast *enough*, simple to use, has top class tooling, and please just
[read the list of std pakages](https://pkg.go.dev/std).

Plus I'm very familiar with it and I've used it a lot.

My biggest problem Go is it's sub-par typesystem.
Generics were a huge step, but we're not all the way there.

# Why Postgres?

I agonized about this for longer than I should.
I **really** wanted to use
[modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite)
which is a library I really like.

But the deciding factor was that Postgres has more features,
is even *more* robust, and I'm already connecting to several other services.

Doubling down on PG early will save me a lot of hassle later.
It can easily be re-used for search, blob storage, caching, etc.

SQL databases like Postgres are giant balls of mutable state.
This is fine, this is good, the DBMS is there to handle all this.
We can leverage that to shove state out of our program.

## Why is the Config in the DB

I consider the configuration of the program to be piece of mutable state.
So like all the other mutable state it belongs in the database.

# Why Passwordless?

Oh how I wish I could go back in time and remake the web
[a la Hypnospace](https://store.steampowered.com/app/844590/Hypnospace_Outlaw/).

Passwords were never a good idea.
They're a poor imitation of proper cryptography practices.

# Dependency Choices

See [`go.mod`](../go.mod)

## Go Standard Library

I'm counting /x/ in here because it's "official" too.

## Gorm

This one I'm still mixed on but the end choice was made because I
***HATE*** migrations systems. No I don't have a better solution.

So instead I rely on the worse solution of Gorm automigrate and writing
a pile of SQL inside Go code. Maybe I'll clean it up someday...

## Readability and Obelisk

Both of these come from the wonderful [Shiori project](https://github.com/go-shiori/shiori)
and I am happy to say that they both work great!

## Unpoly

For you HTMX fans, try out Unpoly.
At under 50kb from a cached CDN, it's pretty small.
It's a wonderful progressive-enhancement framework that lets me easily build
simple and snappy websites.

Unpoly is loaded from JsDelivr