package diff

import (
	"bytes"
	"testing"
)

func TestDiff(t *testing.T) {
	cmd := NewDiffCommand()
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"--src", "../../..", "--patchRaw", `
diff --git a/cmd/sibyl/root.go b/cmd/sibyl/root.go
index 3bf552b..339495a 100644
--- a/cmd/sibyl/root.go
+++ b/cmd/sibyl/root.go
@@ -1,9 +1,9 @@
 package main
 
 import (
-	"fmt"
-	"github.com/spf13/cobra"
 	"log"
+
+	"github.com/spf13/cobra"
 )
 
 var rootCmd = &cobra.Command{
@@ -11,7 +11,7 @@
 var rootCmd = &cobra.Command{
 	Short: "sibyl cmd",
 	Long:  "sibyl cmd",
 	Run: func(cmd *cobra.Command, args []string) {
-		fmt.Println("Root cmd from sibyl 2")
+		cmd.Help()
 	},
 }

`})
	cmd.Execute()
}
