package main

import (
	"fmt"

	"github.com/ddddami/bindle/strutil"
)

func main() {
	// URL Sanitization
	originalTitle := "This is a Test     Article!--- With Spaces & Special Characters?"
	sanitizedURL := strutil.SanitizeURL(originalTitle)
	fmt.Printf("Original: %s\nSanitized URL: %s\n\n", originalTitle, sanitizedURL)

	// Truncate string
	longText := "This is a very long text that needs to be truncated because it exceeds the maximum allowed length."
	truncated := strutil.Truncate(longText, 30)
	fmt.Printf("Original: %s\nTruncated: %s\n\n", longText, truncated)

	// Strip HTML
	htmlText := "<p>This is <strong>HTML</strong> content with <a href='#'>a link</a> and &quot;entities&quot;.</p>"
	plainContent := strutil.StripHTML(htmlText)
	fmt.Printf("Original: %s\nStripped HTML: %s\n\n", htmlText, plainContent)

	// Format filename
	badFilename := "User's Document: Important! (2023).txt"
	goodFilename := strutil.FormatFilename(badFilename)
	fmt.Printf("Original: %s\nSafe filename: %s\n\n", badFilename, goodFilename)

	// Email validation
	validEmails := []string{"user@example.com", "name.surname@domain.co.uk"}
	invalidEmails := []string{"not-an-email", "missing@domain", "@nocuser.com"}

	fmt.Println("Valid emails:")
	for _, email := range validEmails {
		fmt.Printf("  %s: %v\n", email, strutil.IsValidEmail(email))
	}

	fmt.Println("Invalid emails:")
	for _, email := range invalidEmails {
		fmt.Printf("  %s: %v\n", email, strutil.IsValidEmail(email))
	}
}
