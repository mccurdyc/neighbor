# FAQ
---

## My script requires `sudo` access during execution, what should I do?

This command can be a script that you have written, but keep in mind that if the script
needs `sudo` access, you should add it to the `sudoers` file because requiring you
to type your password for each project is blocking.

To add a script to the `sudoers` file add a line similar to the following:
```bash
your_username ALL=NOPASSWD:/absolute/path/to/script
```

Then, where you will invoke neighbor from, you will need to `source` your `sudoers`
file to acquire the recent changes.


