package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
)

var (
	ERR_PARSE_BODY_EMPTY                          = fmt.Errorf("Parse Body Empty")
	ERR_RESULT_ALREADYEXISTS                      = fmt.Errorf("Result Line Already Exists")
	ERR_RESULT_DOESNTEXIST                        = fmt.Errorf("Result Object Doesn't Exist")
	ERR_RESULT_STILLEXISTS                        = fmt.Errorf("Result Object Still Exists")
	ERR_RESULTS_NOTFOUND                          = fmt.Errorf("Results NOT Found!!!")
	ERR_PARSE_RESULT_START_ID_START_NOTFOUND      = fmt.Errorf("Parse Result_Start ID Start NOT Found!")
	ERR_PARSE_RESULT_START_ID_END_NOTFOUND        = fmt.Errorf("Parse Result_Start ID End NOT Found!")
	ERR_PARSE_RESULT_START_ID_INVALID             = fmt.Errorf("Parse Result_Start ID Invalid!")
	ERR_PARSE_RESULT_START_URL_NOTFOUND           = fmt.Errorf("Parse Result_Start URL NOT Found!")
	ERR_PARSE_RESULT_START_TITLE_START_NOTFOUND   = fmt.Errorf("Parse Result_Start Title Start NOT Found!")
	ERR_PARSE_RESULT_START_TITLE_END_NOTFOUND     = fmt.Errorf("Parse Result_Start Title End NOT Found!")
	ERR_PARSE_RESULT_END_START_NOTFOUND           = fmt.Errorf("Parse Result_End ID Start NOT Found!")
	ERR_PARSE_RESULT_END_TITLE_START_NOTFOUND     = fmt.Errorf("Parse Result_End Title Start NOT Found!")
	ERR_PARSE_RESULT_END_USER_START_NOTFOUND      = fmt.Errorf("Parse Result_End User Start NOT Found!")
	ERR_PARSE_RESULT_END_POINTS_INVALID           = fmt.Errorf("Parse Result_End Points Invalid!")
	ERR_PARSE_RESULT_END_USER_NOTFOUND            = fmt.Errorf("Parse Result_End User NOT Found!")
	ERR_PARSE_RESULT_END_AGE_PREFIX_NOTFOUND      = fmt.Errorf("Parse Result_End Age Prefix NOT Found!")
	ERR_PARSE_RESULT_END_AGE_SUFFIX_NOTFOUND      = fmt.Errorf("Parse Result_End Age Suffix NOT Found!")
	ERR_PARSE_RESULT_END_AGE_START_NOTFOUND       = fmt.Errorf("Parse Result_End Age Start NOT Found!")
	ERR_PARSE_RESULT_END_AGE_INVALID              = fmt.Errorf("Parse Result_End Age Invalid!")
	ERR_PARSE_RESULT_END_COMMENTS_PREFIX_NOTFOUND = fmt.Errorf("Parse Result_End Comments Prefix NOT Found!")
	ERR_PARSE_RESULT_END_COMMENTS_START_NOTFOUND  = fmt.Errorf("Parse Result_End Comments Start NOT Found!")
	ERR_PARSE_RESULT_END_COMMENTS_NOTFOUND        = fmt.Errorf("Parse Result_End Comments NOT Found!") // we have to either find a suffix or none
	ERR_PARSE_RESULT_END_COMMENTS_INVALID         = fmt.Errorf("Parse Result_End Comments Invalid!")
)

var (
	parser_result_closing_tag = []byte(`>`)
	parser_result_quote       = []byte(`"`)
	// start
	parser_result_start_prefix          = []byte(`<td align="right" valign="top" class="title"><span class="rank">`)
	parser_result_start_id_upvote_start = []byte(`vote?id=`)
	parser_result_start_id_upvote_end   = []byte(`&amp;how=up&amp;goto=news'><div class='votearrow' title='upvote'></div></a></center></td><td class="title"><a href="`)
	parser_result_start_url_end         = []byte(`" class="storylink"`) // we can't close this, some links end with `rel="nofollow">`
	parser_result_start_title_end       = []byte(`</a>`)                // we can't include span, ask posts don't include it, we also don't need to since there's no more interesting data on this line
	// end
	parser_result_end_prefix      = []byte(`<span class="score" id="score_`)
	parser_result_user_start      = []byte(` points</span> by <a href="user?id=`)
	parser_result_age_prefix      = []byte(`<span class="age"><a href="item?id=`)
	parser_result_age_suffix      = []byte(`</a>`) // we're going to skip to the end of the href, the body can include ` day/days, hour/hours, minute/minutes, second/seconds ago`
	parser_result_comments_prefix = []byte(`hide</a> | <a href="item?id=`)
	parser_result_comments_suffix = []byte(`&nbsp;`)
	parser_result_comments_none   = []byte(`discuss`)
)

func Parse(
	// alternatively we could accept a reader, but a byte slice might be more useful for now
	data []byte,
) (
	[]*Result,
	error,
) {
	if len(data) == 0 {
		return nil, ERR_PARSE_BODY_EMPTY
	}
	results := []*Result{}
	// create a reader
	reader := bufio.NewReader(bytes.NewReader(data))
	// we're going to store a reference to a result here
	// this way as we alternate between lines that the data spans over we can still access this object
	// if the object SHOULD or SHOULDN'T exist we will return an error!
	var result *Result
	// parse content
	for {
		var l []byte // only used for readline
		var buffer bytes.Buffer
		var isPrefix bool
		var err error
		for {
			// ReadLine tries to return a single line, not including the end-of-line bytes.
			// If the line was too long for the buffer then isPrefix is set and the beginning of the line is returned.
			// The rest of the line will be returned from future calls. isPrefix will be false when returning the last fragment of the line.
			// The returned buffer is only valid until the next call to ReadLine.
			// ReadLine either returns a non-nil line or it returns an error, never both.
			//
			// The text returned from ReadLine does not include the line end ("\r\n" or "\n").
			l, isPrefix, err = reader.ReadLine()
			buffer.Write(l)
			// If we've reached the end of the line, stop reading.
			if !isPrefix {
				break
			}
			// If we're just at the EOF, break
			if err != nil {
				break
			}
		}
		if err == io.EOF {
			// done reading
			break
		} else if err != nil {
			// unknown error
			return nil, err
		}
		// we can parse an individual line now
		//
		// I'm going to try to avoid regular expressions since most of the ycombinator html is very simple and fits on a single line or two
		// I'm also going to trim whitespace so individual lines are easier to manage
		line := bytes.TrimSpace(buffer.Bytes())
		if bytes.HasPrefix(line, parser_result_start_prefix) {
			// RESULT STARTING LINE!!!!
			if result.IsValid() {
				// our previous line was a result
				// WE CAN'T HAVE MULTIPLE RESULT LINES!!!
				return nil, ERR_RESULT_ALREADYEXISTS
			}
			// new result object
			result = &Result{}
			// trim prefix
			line = line[len(parser_result_start_prefix):]
			// split ID start
			split := bytes.SplitAfterN(line, parser_result_start_id_upvote_start, 2)
			if len(split) != 2 {
				// ID split start NOT found
				return nil, ERR_PARSE_RESULT_START_ID_START_NOTFOUND
			}
			// found the start of our ID
			// discard the split prefix
			line = split[1]
			// split the end of our ID, until our URL
			split = bytes.SplitAfterN(line, parser_result_start_id_upvote_end, 2)
			if len(split) != 2 {
				// ID split start NOT found
				return nil, ERR_PARSE_RESULT_START_ID_END_NOTFOUND
			}
			// found the end of our ID
			// trim our id split suffix and validate our ID
			id, err := strconv.ParseUint(string(split[0][:len(split[0])-len(parser_result_start_id_upvote_end)]), 10, 64)
			if err != nil {
				// return our own error
				return nil, ERR_PARSE_RESULT_START_ID_INVALID
			}
			// set ID
			result.ID = id
			// parse URL end
			line = split[1]
			split = bytes.SplitAfterN(line, parser_result_start_url_end, 2)
			if len(split) != 2 {
				// URL split start NOT found
				return nil, ERR_PARSE_RESULT_START_URL_NOTFOUND
			}
			// found our URL
			// we're not going to validate our URL, but we could with net/url
			// trim and set URL
			result.URL = string(split[0][:len(split[0])-len(parser_result_start_url_end)])
			// split starting title
			split = bytes.SplitAfterN(line, parser_result_closing_tag, 2)
			if len(split) != 2 {
				// title split start NOT found
				return nil, ERR_PARSE_RESULT_START_TITLE_START_NOTFOUND
			}
			// found the start of our title
			// discard the split prefix
			line = split[1]
			// split our ending title
			split = bytes.SplitAfterN(line, parser_result_start_title_end, 2)
			if len(split) != 2 {
				// title split end NOT found
				return nil, ERR_PARSE_RESULT_START_TITLE_END_NOTFOUND
			}
			// trim and set title
			result.Title = string(split[0][:len(split[0])-len(parser_result_start_title_end)])
			// we can discard the remainder of this line, there's no extra data within it
		} else if bytes.HasPrefix(line, parser_result_end_prefix) {
			// RESULT ENDING LINE!!!!
			if !result.IsValid() {
				// our previous line WAS NOT a result!!!
				// we need a result object to use!!!
				return nil, ERR_RESULT_DOESNTEXIST
			}
			// trim prefix
			line = line[len(parser_result_end_prefix):]
			// split ending tag, points start
			split := bytes.SplitAfterN(line, parser_result_closing_tag, 2)
			if len(split) != 2 {
				// title split start NOT found
				return nil, ERR_PARSE_RESULT_END_TITLE_START_NOTFOUND
			}
			// found the start of our points
			// discard the split prefix
			line = split[1]
			// split the end of our points until we reach our user
			split = bytes.SplitAfterN(line, parser_result_user_start, 2)
			if len(split) != 2 {
				// ID split start NOT found
				return nil, ERR_PARSE_RESULT_END_USER_START_NOTFOUND
			}
			// found our Points
			// trim our points split suffix and validate our Points
			points, err := strconv.ParseUint(string(split[0][:len(split[0])-len(parser_result_user_start)]), 10, 64)
			if err != nil {
				// return our own error
				return nil, ERR_PARSE_RESULT_END_POINTS_INVALID
			}
			// set Points
			result.Points = points
			line = split[1]
			// split ending quote, user ends here
			split = bytes.SplitAfterN(line, parser_result_quote, 2)
			if len(split) != 2 {
				// split ending quote NOT found
				return nil, ERR_PARSE_RESULT_END_USER_NOTFOUND
			}
			// found the end of our user
			// set User
			result.User = string(split[0][:len(split[0])-len(parser_result_quote)])
			line = split[1]
			// split and discard age prefix
			split = bytes.SplitAfterN(line, parser_result_age_prefix, 2)
			if len(split) != 2 {
				// split NOT found
				return nil, ERR_PARSE_RESULT_END_AGE_PREFIX_NOTFOUND
			}
			line = split[1]
			// our userid is repeated, but we're going to ignore it and split on first closing tag
			split = bytes.SplitAfterN(line, parser_result_closing_tag, 2)
			if len(split) != 2 {
				// split NOT found
				return nil, ERR_PARSE_RESULT_END_AGE_START_NOTFOUND
			}
			line = split[1]
			// we're going to split on our age suffix now
			split = bytes.SplitAfterN(line, parser_result_age_suffix, 2)
			if len(split) != 2 {
				// split NOT found
				return nil, ERR_PARSE_RESULT_END_AGE_SUFFIX_NOTFOUND
			}
			// set age
			result.Age = string(bytes.ToLower(split[0][:len(split[0])-len(parser_result_age_suffix)]))
			line = split[1]
			// we have to find the comments now
			// we're going to split on our comments prefix now
			split = bytes.SplitAfterN(line, parser_result_comments_prefix, 2)
			if len(split) != 2 {
				// split NOT found
				return nil, ERR_PARSE_RESULT_END_COMMENTS_PREFIX_NOTFOUND
			}
			line = split[1]
			// our userid is repeated, but we're going to ignore it and split on first closing tag
			split = bytes.SplitAfterN(line, parser_result_closing_tag, 2)
			if len(split) != 2 {
				// split NOT found
				return nil, ERR_PARSE_RESULT_END_COMMENTS_START_NOTFOUND
			}
			line = split[1]
			// we'll either have a prefix of `discuss` here for no comments
			// or we'll have to split on &amp to get the amount of comments
			if bytes.HasPrefix(line, parser_result_comments_none) {
				// no comments!
				// this is fine
			} else {
				// split comments suffix
				split = bytes.SplitAfterN(line, parser_result_comments_suffix, 2)
				if len(split) != 2 {
					// split NOT found
					return nil, ERR_PARSE_RESULT_END_COMMENTS_NOTFOUND
				}
				// validate comments
				comments, err := strconv.ParseUint(string(split[0][:len(split[0])-len(parser_result_comments_suffix)]), 10, 64)
				if err != nil {
					// return our own error
					return nil, ERR_PARSE_RESULT_END_COMMENTS_INVALID
				}
				// we can discard the remainder of this line
				result.Comments = comments
			}
			// append result
			results = append(results, result)
			// unset current result
			result = nil
		} else {
			// discard this line
			// this is fine!
		}
	}
	if result.IsValid() {
		// a result object still exists
		// this is unusable
		return nil, ERR_RESULT_STILLEXISTS
	}
	if len(results) == 0 {
		// no results found!
		return nil, ERR_RESULTS_NOTFOUND
	}
	// success!!
	return results, nil
}
