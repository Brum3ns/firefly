package httpdiff

import "github.com/Brum3ns/firefly/pkg/httpnode"

type HTMLNodeDiff struct {
	TagStartHits       int               `json:"tagstarthits"`
	TagEndHits         int               `json:"tagendhits"`
	TagSelfCloseHits   int               `json:"tagselfclosehits"`
	WordsHits          int               `json:"wordshits"`
	CommentHits        int               `json:"commenthits"`
	AttributeHits      int               `json:"attributehits"`
	AttributeValueHits int               `json:"attributevaluehits"`
	HTMLNode           httpnode.HTMLNode `json:"htmlnode"`
}

type HTMLResult struct {
	OK        bool         `json:"ok"`
	Appear    HTMLNodeDiff `json:"appear"`
	Disappear HTMLNodeDiff `json:"disappear"`
}

// Run the [diff]erence enumiration process for the HTML node
func (hdiff *HttpDiff) GetHTMLNodeDiff(htmlNode httpnode.HTMLNode) HTMLResult {
	totalHits := 0
	storage := struct {
		appearHits    int
		disappearHits int
		appear        []diffNode
		disappear     []diffNode
	}{}

	//!Note : Order for "current" and "known" togther with the list length *MUST* be the same:
	current := [7]diffNode{
		{data: htmlNode.TagStart},
		{data: htmlNode.TagEnd},
		{data: htmlNode.TagSelfClose},
		{data: htmlNode.Words, checkRandomness: true},
		{data: htmlNode.Comment, checkRandomness: true},
		{data: htmlNode.Attribute},
		{data: htmlNode.AttributeValue, checkRandomness: true},
	}
	known := [7]map[string][]int{
		hdiff.config.Merge.HTMLMergeNode.TagStart,
		hdiff.config.Merge.HTMLMergeNode.TagEnd,
		hdiff.config.Merge.HTMLMergeNode.TagSelfClose,
		hdiff.config.Merge.HTMLMergeNode.Words,
		hdiff.config.Merge.HTMLMergeNode.Comment,
		hdiff.config.Merge.HTMLMergeNode.Attribute,
		hdiff.config.Merge.HTMLMergeNode.AttributeValue,
	}

	for i := 0; i < len(current); i++ {
		//Detect difference
		diffAppear, diffDisappear := hdiff.nodeDiff(current[i], known[i], hdiff.config.Payload)

		storage.appear = append(storage.appear, diffAppear)
		storage.appearHits += diffAppear.hit

		storage.disappear = append(storage.disappear, diffDisappear)
		storage.disappearHits += diffDisappear.hit

		totalHits += (diffAppear.hit + diffDisappear.hit)
	}

	return HTMLResult{
		// !Note : The order for this *MUST* follow the same as "known" and "current" above
		OK: (totalHits > 0),
		Appear: HTMLNodeDiff{
			TagStartHits:       storage.appear[0].hit,
			TagEndHits:         storage.appear[1].hit,
			TagSelfCloseHits:   storage.appear[2].hit,
			WordsHits:          storage.appear[3].hit,
			CommentHits:        storage.appear[4].hit,
			AttributeHits:      storage.appear[5].hit,
			AttributeValueHits: storage.appear[6].hit,
			HTMLNode: httpnode.HTMLNode{
				TagStart:       storage.appear[0].data,
				TagEnd:         storage.appear[1].data,
				TagSelfClose:   storage.appear[2].data,
				Words:          storage.appear[3].data,
				Comment:        storage.appear[4].data,
				Attribute:      storage.appear[5].data,
				AttributeValue: storage.appear[6].data,
			},
		},
		Disappear: HTMLNodeDiff{
			TagStartHits:       storage.disappear[0].hit,
			TagEndHits:         storage.disappear[1].hit,
			TagSelfCloseHits:   storage.disappear[2].hit,
			WordsHits:          storage.disappear[3].hit,
			CommentHits:        storage.disappear[4].hit,
			AttributeHits:      storage.disappear[5].hit,
			AttributeValueHits: storage.disappear[6].hit,
			HTMLNode: httpnode.HTMLNode{
				TagStart:       storage.disappear[0].data,
				TagEnd:         storage.disappear[1].data,
				TagSelfClose:   storage.disappear[2].data,
				Words:          storage.disappear[3].data,
				Comment:        storage.disappear[4].data,
				Attribute:      storage.disappear[5].data,
				AttributeValue: storage.disappear[6].data,
			},
		},
	}
}
