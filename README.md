This is a very simple golang program to upload a small file into an Azure Storage Account. 
It uses a Device Token and Azure Active Directory to authenticate the user.  The following
are required on the command line.

This program requires an Application Registration to allow the program to impersonate a user.
The enclosed bash script should give you a good start but is not guaranteed to work in your 
environment as it could require Admin consent from the Tenant owner.

```
Usage of simpleaad:
  -appid string
        App Registration Id (Required)
  -container string
        Container for Upload (Required)
  -store string
        Storage Acct for Upload (Required)
  -subid string
        Azure SubscriptionId (Required)
  -tenantid string
        Tenant Id (Required)
```
