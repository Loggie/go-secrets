# go-secrets

Small Go client for org.freedesktop.secrets (Secret Service API over DBus).

## Requirements

- Linux desktop/session with Secret Service available
- Session DBus running
- A secret service backend like gnome-keyring or KWallet

## Install

    go get github.com/Loggie/go-secrets

## Example: Add Secret To A Collection

    item, err := secrets.Add(
        "default",
        "Demo DB Password",
        "text/plain",
        "super-secret-value",
        secrets.Attributes{
            "app":  "demo-cli",
            "name": "db-password",
        },
    )

    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("added secret:", item.Path())

## Example: Search For Secret Using Attributes

    value, err := secrets.Get(secrets.Attributes{
        "app":  "demo-cli",
        "name": "db-password",
    })

    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("retrieved secret:", value)


## Notes

- Use stable attributes. Retrieval by helper depends on matching attributes.
- If multiple items match, helper returns first one.
- `Add("", ...)` uses the default collection.
