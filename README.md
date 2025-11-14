# golang-http-invalid-context-length

## Investigation of Go HTTP Library Content-Length Behavior

This repository investigates whether the Go HTTP library allows producing invalid HTTP responses where the body is longer than the Content-Length header indicates.

### Test Cases

All test cases use a DELETE handler that writes "helloworld" (10 bytes) as the response body:

- **Test a)** No content-length, Status 200
- **Test b)** No content-length, Status 204
- **Test c)** Content-Length: 2, Status 204
- **Test d)** Content-Length: 11, Status 204
- **Test e)** Content-Length: 2, Status 200

### Key Findings

Running `go test -v` produces detailed output showing:

1. **Status 200 with no Content-Length (Test a)**: Go automatically adds `Content-Length: 10` and sends the full body.

2. **Status 204 with no Content-Length (Test b)**: Go strips the body completely (204 No Content must not have a body per RFC). The handler's Write() call is ignored.

3. **Status 204 with Content-Length: 2 (Test c)**: Go ignores the Content-Length header and strips the body (204 behavior takes precedence).

4. **Status 204 with Content-Length: 11 (Test d)**: Go ignores the Content-Length header and strips the body (204 behavior takes precedence).

5. **Status 200 with Content-Length: 2 (Test e)**: Go sends only 2 bytes but tries to write 10 bytes. The client receives an `unexpected EOF` error when trying to read the body, as only 2 bytes are available but 10 were written by the handler.

### Conclusion

The Go HTTP library has built-in protections:
- It automatically calculates Content-Length when not provided for status codes that allow bodies
- It enforces RFC-compliant behavior for 204 No Content (strips body regardless of Content-Length)
- When Content-Length is manually set lower than the actual body, it truncates the response, causing client-side errors

However, **it is possible to create invalid responses** (Test e) where the server tries to write more data than Content-Length indicates, resulting in truncated responses and client errors.

### Running Tests

```bash
go test -v
```

Or use GitHub Actions to see the results in CI.