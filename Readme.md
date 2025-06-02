# ha runner

Incredibly mininal remote home assistant command runner. Currently just what I need. Might dev more if anyone else finds this.

Basic premise is, listens for entity events and then runs commands. Great for a push button helper that does something on a computer somewhere!

See the example config TOML for details, expects config file as it's argument.

Expects a long lived token for your home assistant in the env variable `HA_AUTH_TOKEN`.
