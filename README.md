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

1. **Status 200 with no Content-Length (Test a)**: Go automatically adds `Content-Length: 10` and sends the full body. ✓ PASS

2. **Status 204 with no Content-Length (Test b)**: Go strips the body completely (204 No Content must not have a body per RFC). The handler's Write() call is ignored. ✓ PASS - body is empty

3. **Status 204 with Content-Length: 2 (Test c)**: Go ignores the Content-Length header and strips the body (204 behavior takes precedence). ✓ PASS - body is empty

4. **Status 204 with Content-Length: 11 (Test d)**: Go ignores the Content-Length header and strips the body (204 behavior takes precedence). ✓ PASS - body is empty

5. **Status 200 with Content-Length: 2 (Test e)**: Go truncates the response to match the Content-Length header. The client receives an `unexpected EOF` error when trying to read, as it receives 0 bytes instead of the promised 2. ✓ PASS - body never exceeds Content-Length

### Test Assertions

Each test validates the following invariants:
- **204 No Content responses MUST have empty body** - Tests will FAIL if a 204 response contains any body data
- **Body length MUST NOT exceed Content-Length header** - Tests will FAIL if the received body is larger than the Content-Length header specifies

### Conclusion

The Go HTTP library enforces HTTP protocol constraints:
- It automatically calculates Content-Length when not provided for status codes that allow bodies
- It enforces RFC-compliant behavior for 204 No Content (strips body regardless of Content-Length)
- It respects manually-set Content-Length headers by truncating responses (never sends more than Content-Length specifies)
- When Content-Length is manually set lower than the actual body, it truncates the response, causing client-side errors

**All tests pass**, confirming that Go's HTTP library never sends a body larger than Content-Length indicates, and always returns an empty body for 204 responses.

### Running Tests

```bash
go test -v
```

Or use GitHub Actions to see the results in CI.