package upload

import (
	"encoding/xml"
	"fmt"
	"os"
	"os/exec"
)

type Project struct {
	XMLName    xml.Name `xml:"project"`
	GroupID    string   `xml:"groupId"`
	ArtifactID string   `xml:"artifactId"`
	Version    string   `xml:"version"`
}

type MvnUploadJars struct{}

func (MvnUploadJars) Upload(jarPath, pomPath, repoURL, repoID string) error {
	// 读取 POM 文件
	pomContents, err := os.ReadFile(pomPath)
	if err != nil {
		return fmt.Errorf("could not read POM file: %w", err)
	}

	// 解析 POM 文件
	var project Project
	if err := xml.Unmarshal(pomContents, &project); err != nil {
		return fmt.Errorf("could not unmarshal POM file: %w", err)
	}

	// 构建 Maven 命令
	var cmd *exec.Cmd
	if jarPath != "" {
		cmd = exec.Command("mvn", "deploy:deploy-file",
			"-DgroupId="+project.GroupID,
			"-DartifactId="+project.ArtifactID,
			"-Dversion="+project.Version,
			"-Dpackaging=jar",
			"-Dfile="+jarPath,
			"-DpomFile="+pomPath,
			"-Durl="+repoURL,
			"-DrepositoryId="+repoID,
		)
	} else if pomPath != "" {
		cmd = exec.Command("mvn", "deploy:deploy-file",
			"-DgroupId="+project.GroupID,
			"-DartifactId="+project.ArtifactID,
			"-Dversion="+project.Version,
			"-Dpackaging=pom",
			"-Dfile="+pomPath,
			"-Durl="+repoURL,
			"-DrepositoryId="+repoID,
		)
	}

	// 执行 Maven 命令
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to deploy %s: %w\nOutput: %s", jarPath+":"+pomPath, err, output)
	}
	fmt.Printf("Successfully deployed: %s\n", jarPath+":"+pomPath)
	return nil
}
