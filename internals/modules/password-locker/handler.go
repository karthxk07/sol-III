package passwordlocker

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {

	//bind route /google to the main engine
	passwordlocker_router := r.Group("/passwordlocker")

	passwordlocker_router.GET("/gen", func(c *gin.Context) {
		// 1. Resolve path relative to binary
		exePath := "/home/karthik/Projects/sol-III/internals/modules/password-locker/"
		targetPath := filepath.Join(filepath.Dir(exePath), ".encoded_output")

		// 2. Check for file existence
		if _, err := os.Stat(targetPath); err == nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  "exists",
				"message": "encoded file already exists",
			})
			return
		}

		// 3. Generate hex-32 number (16 bytes = 32 hex characters)
		bytes := make([]byte, 16)
		if _, err := rand.Read(bytes); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate random sequence"})
			return
		}
		hexString := hex.EncodeToString(bytes)

		// 4. Encode and save
		if err := SaveEncodedPayload(hexString); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 5. Success response
		c.JSON(http.StatusOK, gin.H{
			"status":    "generated",
			"generated": hexString,
		})
	})

	passwordlocker_router.GET("/get", func(c *gin.Context) {
		// 1. Check Indian Standard Time
		loc, err := time.LoadLocation("Asia/Kolkata")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Timezone lookup failed"})
			return
		}

		now := time.Now().In(loc)
		day := now.Weekday()
		hour := now.Hour()

		if (day == time.Friday && hour <= 18) && day != time.Saturday && day != time.Sunday {
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "denied",
				"message": "It's a weekday. Get back to work; the vault only opens when the world rests.",
			})
			return
		}

		// 2. Resolve path for the hidden file
		exePath := "/home/karthik/Projects/sol-III/internals/modules/password-locker/"
		targetPath := filepath.Join(filepath.Dir(exePath), ".encoded_output")

		fmt.Println(os.Stat(targetPath))
		// 3. Check if the hidden file exists
		if _, err := os.Stat(targetPath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"status":      "missing_data",
				"message":     "The password file is non-existent, much like your preparation for this request.",
				"instruction": "Go generate a new password and reset it at your end before trying to knock on this door again.",
			})
			return
		}

		// 4. Continue to decode and return
		decodedString, err := ReadAndDecodePayload() // Ensure this helper points to ".encoded_output"
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decrypt the void"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":   "success",
			"password": decodedString,
		})
	})

	passwordlocker_router.GET("/iamloser", func(c *gin.Context) {
		// 1. Get the query parameter
		lostDate := c.Query("iloston")
		if lostDate == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "You must admit when you lost. Missing 'iloston' param."})
			return
		}

		// Resolve path relative to binary
		exePath := "/home/karthik/Projects/sol-III/internals/modules/password-locker/"
		lostFilePath := filepath.Join(filepath.Dir(exePath), "lost.txt")

		// 2. Read existing content (Ignore error if file doesn't exist yet)
		existingContent, _ := os.ReadFile(lostFilePath)

		// 3. Prepare the updated log
		newLine := fmt.Sprintf("you lost on %s", lostDate)
		var updatedContent string
		if len(existingContent) > 0 {
			updatedContent = string(existingContent) + "; " + newLine
		} else {
			updatedContent = newLine
		}

		// 4. Run the decode function to get the current password
		password, err := ReadAndDecodePayload()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve password", "detail": err.Error()})
			return
		}

		// 5. Save the updated log back to the file
		if err := os.WriteFile(lostFilePath, []byte(updatedContent), 0644); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log your loss"})
			return
		}

		// 6. Return the password and the full loss history
		c.JSON(http.StatusOK, gin.H{
			"status":   "admitted",
			"password": password,
			"history":  updatedContent,
		})
	})
}
