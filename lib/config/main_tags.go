package config

// main tags is a config, load from database

var mainTags []string

var mainTagSet map[string]bool

func SetMainTags(tags []string) {
	mainTags = []string{}
	mainTagSet := map[string]bool{}
	for _, t := range tags {
		if !mainTagSet[t] {
			mainTags = append(mainTags, t)
			mainTagSet[t] = true
		}
	}
}

func GetMainTags() []string {
	return mainTags
}

func MainTags() []string {
	return mainTags
}

func SplitTags(tags ...string) (mains []string, subs []string) {
	repeat := map[string]bool{}
	for _, tag := range tags {
		if repeat[tag] {
			continue
		}
		repeat[tag] = true

		if mainTagSet[tag] {
			mains = append(mains, tag)
		} else {
			subs = append(subs, tag)
		}
	}
	return
}
