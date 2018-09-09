package main

import (
	"fmt"
	"os"
)

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../..")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("error loading config file: %+v\n", err)
		os.Exit(1)
	}

	// client := github.NewClient(nil)
	//
	// // res, _, err := client.Search.Repositories(context.Background(), "go test ./ in:readme starts:>=10000 language:go", nil)
	// res, _, err := client.Search.Repositories(context.Background(), "stars:>=10000 language:python", nil)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	//
	// clonedRepos := map[string]string{}
	// fmt.Println(res)
	//
	// // for _, r := range res.Repositories {
	// // 	name := *r.FullName
	// // 	if _, ok := clonedRepos[name]; ok {
	// // 		continue
	// // 	}
	//
	// // p := fmt.Sprintf("github.com/%s", name)
	// // getCmd := exec.Command("go", "get", "-v", "-d", fmt.Sprintf("%s/...", p))
	// // out, err := getCmd.CombinedOutput()
	// // fmt.Printf("combined out:\n%s\n", string(out))
	// // if err != nil {
	// // 	log.Fatalf("cmd.Run() failed with %s\n", err)
	// // }
	// // fmt.Printf("go get %s\n", p)
	// //
	// // clonedPath := fmt.Sprintf("%s/src/%s", os.Getenv("GOPATH"), p)
	// // if err := os.Chdir(clonedPath); err != nil {
	// // 	fmt.Println(err)
	// // 	continue
	// // }
	// // fmt.Printf("changed directory to %s\n", clonedPath)
	// //
	// // cmd := exec.Command("go", "test", "-v", "./...")
	// // out, err = cmd.CombinedOutput()
	// // fmt.Printf("combined out:\n%s\n", string(out))
	// // if err != nil {
	// // 	log.Fatalf("cmd.Run() failed with %s\n", err)
	// // }
	// // fmt.Printf("go test %s\n", p)
	// //
	// // cmd = exec.Command("go", "test", "-covermode=count", "-coverprofile=/tmp/count.out", "-v", "./...")
	// // out, err = cmd.CombinedOutput()
	// // fmt.Printf("combined out:\n%s\n", string(out))
	// // if err != nil {
	// // 	log.Fatalf("cmd.Run() failed with %s\n", err)
	// // }
	// // fmt.Printf("go test %s\n", p)
	//
	// // clonedRepos[name] = p
	//
	// // FIXME: just for now
	// // 	break
	// // }
}
