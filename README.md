## Image preview
A sample REST application written in Go.

### Run
Deploy application into a Docker container:

```make run-dev-docker```

The application's only endpoint is http://\<host\>:9090/api/images/previews

It accepts JSON and multipart form requests.

Sample JSON request:

```json
{
	"images": ["https://avtoreliz.com/wp-content/uploads/2015/04/infiniti-vision-gt-1.jpg"]
}
```

An image is either URL or Base64 encoded image data.

Form request is intended to contain single or repeated `image` text field
with the same possible values as a JSON image.

`image` file uploading is also supported.

### Test
Run tests:

```make test```
