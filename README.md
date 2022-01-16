How to create credentials.json file (google-drive)

1. go to https://console.cloud.google.com with your gmail account
2. choose navigation menu on the left -> APIs and Services -> Credentials -> + Create Credentials
-> OAuth client ID -> Application Type: Desktop App -> Rename then click "Create"
3. your credential will appear in the table list of "OAuth 2.0 Client IDs"
4. click download button of your credential in the table
5. choose "Download JSON"
6. Rename the downloaded file to "credentials.json"

read more about the tutorial here:
+ https://developers.google.com/drive/api/v3/quickstart/go
+ https://developers.google.com/workspace/guides/create-credentials

_______________________________________________________________

How to create credentials.json file (google-storage)

follow the tutorial: https://developers.google.com/workspace/guides/create-credentials#service-account

*note: choose credentials type "Service-Account" in the above document 

_______________________________________________________________

How to config dropbox cloud to connect

dropbox sdk to use in code: https://github.com/dropbox/dropbox-sdk-go-unofficial

1. Sign in dropbox app console, via link: https://www.dropbox.com/developers, you can choose to sign in with email or facebook, it will best if you sign in with the email that already has a dropbox workspace
2. create new app in "App Console"
   1. choose an API -> Scoped access
   2. choose the type access you need -> full dropbox
   3. name your app
3. In tab "settings" of your app
   1. section OAuth2
      1. access token expiration -> no expiration
      2. allow public clients -> allow
4. In tab "permissions" of your app
   1. Individual Scopes -> tick "write" permission for all
5. get access token: In tab settings of your app, OAuth2 section -> Generated access token -> Generate

_______________________________________________________________

How to create azure application to connect and how to get one drive access token

1. config your app in microsoft azure portal, follow the tutorial: https://docs.microsoft.com/en-us/onedrive/developer/rest-api/getting-started/graph-oauth?view=odsp-graph-online.
but here are some steps that you should do in order for the app to work properly
   1. register your app, link: https://aka.ms/AppRegistrations
      1. "Supported account types" should be "Accounts in any organizational directory (Any Azure AD directory - Multitenant) and personal Microsoft accounts (e.g. Skype, Xbox)"
      2. After successfully create your app, you must do these following steps to config your account
         1. **set client credentials**. go to "Certificates and secrets" in the left menu -> choose tab "client secrets" -> new client secret -> fill client secret info -> save
         2. **create a redirect URIs**. it should be a spa (single page application). go to "Authentication" in the left menu -> Add Platform -> single page application -> quick start -> choose the language you want then do exactly the tutorial says
         3. **platform configuration**. go to "Authentication" in the left menu -> in the "Implicit grant and hybrid flows" section, tick "Access tokens" and "ID tokens" -> Advanced settings, "yes" for "live sdk support", "yes" for "Allow public client flows" -> save
         4. **config application scopes**.
            1. create your app scope.  go to "Expose an API" in the left menu -> create 2 scopes 
               1. ".../Files.ReadWrite" for admin and users
               2. ".../offline_access" for admin and users
               3.  ".../Sites.ReadWrite.All" for admin and users
               4. ".../Files.ReadWrite.All" for admin and users, you can follow this document to create your app's scope: https://docs.microsoft.com/en-us/azure/active-directory/develop/quickstart-configure-app-expose-web-apis
            2. Configured permissions. go to "API Permissions" in the left menu
               1. Your app permission. "Add permission" -> My APIS -> choose your app -> choose all scopes that you just created in step 1
               2. Microsoft Graph permission. "Add permission" -> Microsoft APIS -> Microsoft Graph 
                  1. Delegated Permissions -> choose these permissions: email, offline_access, openid, profile, Files.ReadWrite, Files.ReadWrite.All, User.Read, Sites.ReadWrite.All
                  2. Application Permissions -> choose these permissions: Files.ReadWrite.All, Sites.ReadWrite.All
   . Once you are done, your overview should look like this
   ![img.png](img.png)
2. Sign your user in with the specified scopes using the token flow (there is another method "code flow" but we use "token flow" for example because it's easier)
   follow this document: "token flow" section, https://docs.microsoft.com/en-us/onedrive/developer/rest-api/getting-started/graph-oauth?view=odsp-graph-online
   1. ![img_1.png](img_1.png)
      paste this URL in your "postman" app, then fill all the parameters, it should look like this
      ![img_2.png](img_2.png)
      1. copy the url in your postman then paste it into your browser to redirect to get token
         (**Remember to launch your spa first, in step 1.i.b.b**)
      2. after successfully redirect, the token you need will appear in the redirect URL, copy then paste it in the config of your application