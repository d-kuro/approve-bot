package approve

/*
func Approve() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ""},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	files, _, err := client.PullRequests.ListFiles(ctx, "classmethod", "metropolis-kubernetes", 168, &github.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		fmt.Println(*f.Filename)
	}

	bs := []byte("https://github.com/classmethod/metropolis-kubernetes/pull/1")
	reg := regexp.MustCompile("https://github.com(.*)/(.*)/pull/(\\d*)")
	group := reg.FindSubmatch(bs)
	fmt.Println(group)

}
*/
