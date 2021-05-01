# csrf

[![release](https://github.com/idthings/csrf/actions/workflows/release.yaml/badge.svg)](https://github.com/idthings/csrf/actions/workflows/release.yaml)

* Suitable for distributed deployments
* Simple to understand, and implement
* Easy secret rotation
* No datastore required

#### Approach
This is a, kind of, reverse HMAC header.
Using a random token as the 'data'.

The server side generates a random token, which is sent as a hidden form field.
The server uses its secret to create a SHA256 hash of the random token, as a second hidden form field.
When the form is submitted, the server re-computes the hash from the submitted form and compares.

The server secrets can be rotated, regularly, without breaking in-flight form requests.
Server secrets are picked up from the env var CSRF_SALTS, a list in csv format.
Inject these into the Docker container, or the local env.
```
$ echo $CSRF_SALTS
veClIlI6DN5qBeFifxodUt08PEAFXDvb,YuN8WJPWyW8MBG9csENe0SCvMYmKIfzH,0PS0eLn9TqF093fR3pQ4gow9giMLperw
```
This makes it simple to sync any number of secret rotations across many containers/instances of your web app.
Deploy a cron job to all containers ahead of time, which updates the env var.

We always generate new tokens using the left-most secret.
And we walk left to right when validating incoming hashes.
Logging which secret (index) was used, means you can remove right-most secrets when they are no longer in use.
#### Use
```
$ go get github.com/idthings/csrf
$ export CSRF_SALTS='myreallylongsaltthatshouldberandom'
```
Once imported, generating a token:
```
    csrfToken, csrfHash, err := csrf.Generate()
    if err != nil {
        // do something
    }
    
    var templateData = struct{
        CSRFToken string
        CSRFHash string
    }{
        CSRFToken: csrfToken,
        CSRFHash: csrfHash,
    }
```
And using the token/hash pair in a Go html template:
```
<input type="hidden" name="_token" value="{{ .CSRFToken }}" />
<input type="hidden" name="_hash" value="{{ .CSRFHash }}" />
```
Here's a snippet example of usage in a Go web handler, when validating form input requests.
Capturing and logging the keyUsed to validate, means you can keep track of when it's safe to retire old keys.
```
func FormRequestHandler(w http.ResponseWriter, r *http.Request) {

    token := r.FormValue("_token")
    hash := r.FormValue("_hash")
    valid, keyUsed := csrf.Validate(token, hash)
    if valid && keyUsed > 0 {
        log.Info("CSRF: valid use of key", keyUsed)
    }
}
```
