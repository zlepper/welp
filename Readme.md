# Welp
The Ultra Super Mega Simple Feedback system. 

Welp is super simple to install, simple to integrate with, and even easier to use. 

Welp can be simplified to a system that takes a contact address, a message, and some attached files, and saves them 
somewhere for retrieval later on. 
Then it adds some support for a couple of simple use-cases on top, like responding to the user who provided the 
feedback from the build in UI. 


## What welp doesn't do
Any sort of user tracking. If you want tracking information on how your users are actually using your applications,
use something like Google Analytics, or Microsoft Application Insights. 

## How to install
Go to the [releases page on GitHub](https://github.com/zlepper/welp/releases) and download a release for your
server distribution. The releases should also work on any desktop computer, so it's possible to test integrations 
locally without needing a full server.

Once you have a binary downloaded, simply invoke it from the command line to start with the simplest of all
configurations. If no flags are specified, the server will attempt to start on port 8080 on http. 

The following flags are available (can also be listed by passing `--help` to the binary):
|Flag name|Description|Default value|Recommendation|
|---------|-----------|-------------|--------------|
|--config|The config file to persist options in|$HOME/.welp.yaml|Leave this alone for now.|
|--databaseFolderPath|Where to save the "database" when using the flat-file database|db|No reason to change this|
|--emailSenderAddress|The address of the email sender when welp is sending out emails|noreply@noreply.com|
Change this to some mail that is actually watched in your org, because people _will_ respond to it, weather you want it or not.|
|--emailSenderName|The name of the owner of the address of the above mentioned flag|no-reply|Change this, if you changed the above|
|-h, --help|Prints the help for the application||Use if you want a reminder about all the flags|
|--port|The port to run Welp on, will be ignore if --useHttps is passsed|8080|Change if you have port conflicts, or actually want to host on port 80 without https.|
|--saveInterval|How often the flatFile storage should save changes (such as new feedback, or user changes). Lower values provides better guarantee that data 
doesn't get lost, but will decrease performance.|5s|No reason to change this, unless it becomes an issue.|
|--sendGridApiKey|An api key for sendGrid (https://sendgrid.com/). If provided, sendGrid will be used for sending emails.||Set if you want to use SendGrid for sending email.|
|--storageFolderPath|Sets the folder welp should storage uploaded files to, if using flat file storage|storage|No reason to change this|
|--tokenDuration|How long a login token should be valid. Shorter times is probably more secure, but longer times makes it easier for users. 
Defaults to a year.|8760h0m0s (one year)|Decrease this for probably better security, increase for better use experience|
|--useHttps|Enable to automatically use https everywhere. Will add http -> https redirect, and enable HSTS. Not compatible 
with the --port flag. Will automatically fetch https certificates using Lets Encrypt (For this reason Welp needs to be available 
to the general internet)|false|Use this in production, no excuses.|

## General usage
There are two ways to integrate Welp into your other projects, either pop and iframe pointing to the 
`/embed` endpoint of welp. This endpoint returns a small page with the simple feedback inputs, and 
a file uploader. 

Alternatively it is possible to use the API to create a custom design, so it can fit better in the 
overall project. See the following Api documentation section for that. 

## Api documentation
The welp api is very flexible, and accept most types of input, and can return most types of output.
This is all based on what headers are passed along with the requests. If you want to talk json to 
and endpoint, and get xml back (for some reason), pass the `Content-Type: application/json` and
`Accept: application/xml` header along with your request. This works for all most all apis. 
Exceptions to this will be noted for each api call. 
It's possible to send to following kinds of request content types to almost all apis:
`application/json`, `application/xml`, `application/x-www-form-urlencoded`, `multipart/form-data` 
and query parameters. 
The following types of response content can be requested:
`application/json`, `application/xml`, `text/html`. Just be aware that `text/html` is mostly for
browsers, so it might return different response codes, to guide the browser around. 
For programmatic usage, it's recommended to request json or xml. 

### Creating a new feedback entry:
Send a form post to `/`, with the following attributes:
|key|description|
|-----|-----|
|`message`|The message entered by the user|
|`contactAddress`|The contact address, the user wishes to be contacted on. Not required in the api.|
|`files`|The files as multiparts.|

This api does _only_ accept form posts, not json or xml, due to the file upload support. 
It can return all listed types of response. If the request succeeds, 201 status code is returned. 

### Logging in
You need to login to access most of the apis in welp. 
Send a POST requests to `/login` with the following parameters:
|key|description|
|-----|-----|
|`email`|The email of the user|
|`password`|The password of the user|

This api wil always attempt to set a token cookie. 
If the api should return xml or json, the token will be in the request response. 
The token should then be passed back in the `Authorization` header, as `Bearer <token>`. 

### Get feedback list
To get the full list of feedback, send a GET request to `/`. 
This endpoint requires authentication. 

Does not take any parameters. 


## The build the project
Welp can be fully build by simple running 
```
$ go build
```
That will build a binary for your current system. The binary can be invoked from the command line. 

If you have changed the templates, `go generate` needs to be run to update the embedded templates.
If you are working from an IDE like Intellij, I would recommend just running `go generate` as a pre-step for your 
normal run profile, so you can be sure that everything is update to date whenever you start the project. 

