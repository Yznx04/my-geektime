package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/zkep/my-geektime/internal/global"
	"github.com/zkep/my-geektime/internal/model"
	"github.com/zkep/my-geektime/internal/service"
	"github.com/zkep/my-geektime/internal/types/geek"
)

type CourseSummary struct {
	CourseID     string `json:"course_id"`
	ProductID    string `json:"product_id,omitempty"`
	Title        string `json:"title"`
	Subtitle     string `json:"subtitle,omitempty"`
	Author       string `json:"author,omitempty"`
	AuthorIntro  string `json:"author_intro,omitempty"`
	Cover        string `json:"cover,omitempty"`
	ArticleCount int    `json:"article_count"`
	IsVideo      bool   `json:"is_video,omitempty"`
	IsAudio      bool   `json:"is_audio,omitempty"`
	UpdatedAt    int64  `json:"updated_at,omitempty"`
}

type ListCoursesResult struct {
	Total    int64           `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"page_size"`
	Courses  []CourseSummary `json:"courses"`
}

type Chapter struct {
	Index        int    `json:"index"`
	ArticleID    string `json:"article_id"`
	Aid          string `json:"aid,omitempty"`
	Title        string `json:"title"`
	ChapterTitle string `json:"chapter_title,omitempty"`
	Summary      string `json:"summary,omitempty"`
	IsVideo      bool   `json:"is_video,omitempty"`
	VideoTime    string `json:"video_time,omitempty"`
	AudioTime    string `json:"audio_time,omitempty"`
	ReadingTime  int    `json:"reading_time,omitempty"`
	HasContent   bool   `json:"has_content"`
}

type CourseDetail struct {
	CourseSummary
	Intro       string    `json:"intro,omitempty"`
	IntroHTML   string    `json:"intro_html,omitempty"`
	Chapters    []Chapter `json:"chapters"`
	ChapterHint string    `json:"chapter_hint,omitempty"`
}

type ArticleDetail struct {
	ArticleID     string `json:"article_id"`
	Aid           string `json:"aid,omitempty"`
	CourseID      string `json:"course_id,omitempty"`
	CourseTitle   string `json:"course_title,omitempty"`
	Index         int    `json:"index,omitempty"`
	Title         string `json:"title"`
	Subtitle      string `json:"subtitle,omitempty"`
	ChapterTitle  string `json:"chapter_title,omitempty"`
	Author        string `json:"author,omitempty"`
	Summary       string `json:"summary,omitempty"`
	Content       string `json:"content,omitempty"`
	ContentFormat string `json:"content_format,omitempty"`
	ContentChars  int    `json:"content_chars"`
	Truncated     bool   `json:"truncated,omitempty"`
	IsVideo       bool   `json:"is_video,omitempty"`
	VideoTime     string `json:"video_time,omitempty"`
	AudioTime     string `json:"audio_time,omitempty"`
	ReadingTime   int    `json:"reading_time,omitempty"`
	CommentsCount int64  `json:"comments_count"`
}

type CommentReply struct {
	ID        int64  `json:"id"`
	UserName  string `json:"user_name"`
	Content   string `json:"content"`
	LikeCount int64  `json:"like_count"`
	Ctime     int64  `json:"ctime,omitempty"`
	Time      string `json:"time,omitempty"`
}

type CommentItem struct {
	CommentID       int64          `json:"comment_id"`
	UserName        string         `json:"user_name"`
	Content         string         `json:"content"`
	LikeCount       int64          `json:"like_count"`
	DiscussionCount int64          `json:"discussion_count"`
	Ctime           int64          `json:"ctime,omitempty"`
	Time            string         `json:"time,omitempty"`
	Replies         []CommentReply `json:"replies,omitempty"`
}

type CommentsResult struct {
	ArticleID   string        `json:"article_id"`
	Aid         string        `json:"aid"`
	Title       string        `json:"title,omitempty"`
	Total       int64         `json:"total"`
	Page        int           `json:"page"`
	PageSize    int           `json:"page_size"`
	Comments    []CommentItem `json:"comments"`
	HasNext     bool          `json:"has_next,omitempty"`
	RepliesNote string        `json:"replies_note,omitempty"`
}

type SearchMatch struct {
	Type      string `json:"type"`
	ArticleID string `json:"article_id,omitempty"`
	Aid       string `json:"aid,omitempty"`
	Title     string `json:"title,omitempty"`
	Snippet   string `json:"snippet,omitempty"`
}

type SearchCourse struct {
	CourseID string        `json:"course_id"`
	Title    string        `json:"title"`
	Subtitle string        `json:"subtitle,omitempty"`
	Author   string        `json:"author,omitempty"`
	Matches  []SearchMatch `json:"matches"`
	Score    int           `json:"-"`
}

type SearchResult struct {
	Keyword          string         `json:"keyword"`
	TotalCourses     int            `json:"total_courses"`
	TotalArticleHits int            `json:"total_article_hits"`
	TotalCommentHits int            `json:"total_comment_hits"`
	Page             int            `json:"page"`
	PageSize         int            `json:"page_size"`
	Courses          []SearchCourse `json:"courses"`
	SearchScope      string         `json:"search_scope"`
	Hint             string         `json:"hint,omitempty"`
}

func listCourses(ctx context.Context, page, pageSize int) (*ListCoursesResult, error) {
	tx := global.DB.WithContext(ctx).Model(&model.Task{}).
		Where("task_type = ? AND task_pid = ? AND deleted_at = ?", service.TASK_TYPE_PRODUCT, "", 0)
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, err
	}
	var tasks []*model.Task
	if err := tx.Order("updated_at DESC, id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&tasks).Error; err != nil {
		return nil, err
	}
	ret := &ListCoursesResult{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Courses:  make([]CourseSummary, 0, len(tasks)),
	}
	for _, t := range tasks {
		ret.Courses = append(ret.Courses, buildCourseSummary(t))
	}
	return ret, nil
}

func getCourse(ctx context.Context, id string) (*CourseDetail, error) {
	course, err := findCourseTask(ctx, id)
	if err != nil {
		return nil, err
	}
	var product geek.ProductBase
	if len(course.Raw) > 0 {
		_ = json.Unmarshal(course.Raw, &product)
	}
	ret := &CourseDetail{
		CourseSummary: buildCourseSummary(course),
		Intro:         product.Intro,
		IntroHTML:     product.IntroHTML,
		Chapters:      make([]Chapter, 0),
	}
	var tasks []*model.Task
	if err := global.DB.WithContext(ctx).Model(&model.Task{}).
		Where("task_pid = ? AND task_type = ? AND deleted_at = ?", course.TaskId, service.TASK_TYPE_ARTICLE, 0).
		Order("id ASC").
		Find(&tasks).Error; err != nil {
		return nil, err
	}
	for i, t := range tasks {
		var data geek.ArticleData
		if len(t.Raw) > 0 {
			_ = json.Unmarshal(t.Raw, &data)
		}
		info := data.Info
		title := firstNonEmpty(t.TaskName, info.Title)
		ret.Chapters = append(ret.Chapters, Chapter{
			Index:        i + 1,
			ArticleID:    t.TaskId,
			Aid:          articleAid(t, info),
			Title:        title,
			ChapterTitle: info.ChapterTitle,
			Summary:      info.Summary,
			IsVideo:      info.IsVideo,
			VideoTime:    info.Video.Time,
			AudioTime:    info.Audio.Time,
			ReadingTime:  info.ReadingTime,
			HasContent:   len(info.Content) > 0 || len(info.Cshort) > 0,
		})
	}
	ret.ChapterHint = "使用 get_article 并按 chapter 的 article_id 读取正文"
	return ret, nil
}

func getArticle(ctx context.Context, id string, maxChars int) (*ArticleDetail, error) {
	t, err := findArticleTask(ctx, id)
	if err != nil {
		return nil, err
	}
	var data geek.ArticleData
	if len(t.Raw) > 0 {
		if err = json.Unmarshal(t.Raw, &data); err != nil {
			return nil, fmt.Errorf("解析文章数据失败: %w", err)
		}
	}
	info := data.Info
	aid := articleAid(t, info)
	content := info.Content
	if len(content) == 0 {
		content = info.Cshort
	}
	markdown := content
	if len(content) > 0 {
		if converted, err1 := service.HTMLConvertMarkdown(content); err1 == nil {
			markdown = converted
		}
	}
	truncated := false
	if maxChars > 0 {
		markdown, truncated = truncateRunes(markdown, maxChars)
	}
	ret := &ArticleDetail{
		ArticleID:     t.TaskId,
		Aid:           aid,
		CourseID:      t.TaskPid,
		Title:         firstNonEmpty(t.TaskName, info.Title),
		Subtitle:      info.Subtitle,
		ChapterTitle:  info.ChapterTitle,
		Author:        info.Author.Name,
		Summary:       info.Summary,
		Content:       markdown,
		ContentFormat: "markdown",
		ContentChars:  utf8.RuneCountInString(markdown),
		Truncated:     truncated,
		IsVideo:       info.IsVideo,
		VideoTime:     info.Video.Time,
		AudioTime:     info.Audio.Time,
		ReadingTime:   info.ReadingTime,
	}
	if course, err1 := findCourseTask(ctx, t.TaskPid); err1 == nil {
		ret.CourseTitle = buildCourseSummary(course).Title
	}
	var index int64
	global.DB.WithContext(ctx).Model(&model.Task{}).
		Where("task_pid = ? AND task_type = ? AND deleted_at = ? AND id < ?",
			t.TaskPid, service.TASK_TYPE_ARTICLE, 0, t.Id).
		Count(&index)
	ret.Index = int(index) + 1
	if n, err1 := strconv.ParseInt(aid, 10, 64); err1 == nil && n > 0 {
		global.DB.WithContext(ctx).Model(&model.ArticleComment{}).
			Where("aid = ? AND deleted_at = ?", n, 0).
			Count(&ret.CommentsCount)
	}
	return ret, nil
}

func listComments(ctx context.Context, articleRef string, page, pageSize int, includeReplies bool) (*CommentsResult, error) {
	article, err := findArticleTask(ctx, articleRef)
	if err != nil {
		return nil, err
	}
	var data geek.ArticleData
	if len(article.Raw) > 0 {
		_ = json.Unmarshal(article.Raw, &data)
	}
	aid := articleAid(article, data.Info)
	aidInt, err := strconv.ParseInt(aid, 10, 64)
	if err != nil || aidInt <= 0 {
		return nil, fmt.Errorf("文章缺少有效的极客时间 ID: %s", articleRef)
	}
	tx := global.DB.WithContext(ctx).Model(&model.ArticleComment{}).
		Where("aid = ? AND deleted_at = ?", aidInt, 0)
	ret := &CommentsResult{
		ArticleID: article.TaskId,
		Aid:       aid,
		Title:     firstNonEmpty(article.TaskName, data.Info.Title),
		Page:      page,
		PageSize:  pageSize,
		Comments:  make([]CommentItem, 0, pageSize),
	}
	if err = tx.Count(&ret.Total).Error; err != nil {
		return nil, err
	}
	var rows []*model.ArticleComment
	if err = tx.Order("comment_ctime DESC, id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	commentIDs := make([]int64, 0, len(rows))
	for _, row := range rows {
		var c geek.ArticleComment
		if len(row.Raw) > 0 {
			_ = json.Unmarshal(row.Raw, &c)
		}
		item := CommentItem{
			CommentID:       row.Cid,
			UserName:        c.UserName,
			Content:         c.CommentContent,
			LikeCount:       c.LikeCount,
			DiscussionCount: c.DiscussionCount,
			Ctime:           c.CommentCtime,
			Time:            formatUnix(c.CommentCtime),
		}
		for _, r := range c.Replies {
			item.Replies = append(item.Replies, CommentReply{
				ID:       int64(r.ID),
				UserName: firstNonEmpty(r.UserNameReal, r.UserName),
				Content:  r.Content,
				Ctime:    int64(r.Ctime),
				Time:     formatUnix(int64(r.Ctime)),
			})
		}
		if includeReplies && (row.DiscussionCount > 0 || len(item.Replies) > 0) {
			commentIDs = append(commentIDs, row.Cid)
		}
		ret.Comments = append(ret.Comments, item)
	}
	if includeReplies && len(commentIDs) > 0 {
		var discussions []*model.ArticleCommentDiscussion
		if err = global.DB.WithContext(ctx).Model(&model.ArticleCommentDiscussion{}).
			Where("cid IN ? AND deleted_at = ?", commentIDs, 0).
			Order("likes_number DESC, id ASC").
			Find(&discussions).Error; err != nil {
			return nil, err
		}
		grouped := make(map[int64][]CommentReply)
		for _, d := range discussions {
			var raw geek.DiscussionData
			if len(d.Raw) > 0 {
				_ = json.Unmarshal(d.Raw, &raw)
			}
			grouped[d.Cid] = append(grouped[d.Cid], CommentReply{
				ID:        d.Did,
				UserName:  raw.Author.Nickname,
				Content:   raw.Discussion.DiscussionContent,
				LikeCount: d.LikesNumber,
				Ctime:     d.Ctime,
				Time:      formatUnix(d.Ctime),
			})
		}
		for i := range ret.Comments {
			if replies, ok := grouped[ret.Comments[i].CommentID]; ok && len(replies) > 0 {
				ret.Comments[i].Replies = replies
			}
		}
		if len(discussions) > 0 {
			ret.RepliesNote = "replies 来自本地缓存的讨论数据，按点赞数排序"
		}
	}
	if includeReplies && ret.RepliesNote == "" {
		ret.RepliesNote = "本地未缓存到该文章评论的回复数据"
	}
	ret.HasNext = int64(page*pageSize) < ret.Total
	return ret, nil
}

func searchCourses(ctx context.Context, keyword string, page, pageSize int) (*SearchResult, error) {
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return nil, fmt.Errorf("keyword 不能为空")
	}
	like := "%" + strings.ToLower(keyword) + "%"
	lowerKeyword := strings.ToLower(keyword)
	rawExpr := rawTextExpr()

	courseMap := make(map[string]*model.Task)
	var matchedCourses []*model.Task
	if err := global.DB.WithContext(ctx).Model(&model.Task{}).
		Where("task_type = ? AND task_pid = ? AND deleted_at = ?", service.TASK_TYPE_PRODUCT, "", 0).
		Where("LOWER(task_name) LIKE ? OR "+rawExpr+" LIKE ?", like, like).
		Order("id DESC").
		Limit(100).
		Find(&matchedCourses).Error; err != nil {
		return nil, err
	}
	for _, t := range matchedCourses {
		courseMap[t.TaskId] = t
	}

	var titleArticles []*model.Task
	if err := global.DB.WithContext(ctx).Model(&model.Task{}).
		Where("task_type = ? AND deleted_at = ?", service.TASK_TYPE_ARTICLE, 0).
		Where("LOWER(task_name) LIKE ?", like).
		Order("id DESC").
		Limit(100).
		Find(&titleArticles).Error; err != nil {
		return nil, err
	}
	var contentArticles []*model.Task
	if err := global.DB.WithContext(ctx).Model(&model.Task{}).
		Where("task_type = ? AND deleted_at = ?", service.TASK_TYPE_ARTICLE, 0).
		Where("LOWER(task_name) NOT LIKE ?", like).
		Where(rawExpr+" LIKE ?", like).
		Order("id DESC").
		Limit(100).
		Find(&contentArticles).Error; err != nil {
		return nil, err
	}

	var comments []*model.ArticleComment
	if err := global.DB.WithContext(ctx).Model(&model.ArticleComment{}).
		Where("deleted_at = ?", 0).
		Where(rawExpr+" LIKE ?", like).
		Order("comment_ctime DESC").
		Limit(200).
		Find(&comments).Error; err != nil {
		return nil, err
	}
	aids := make([]string, 0, len(comments))
	for _, c := range comments {
		aids = append(aids, strconv.FormatInt(c.Aid, 10))
	}
	articleByAid := make(map[string]*model.Task)
	if len(aids) > 0 {
		var tasks []*model.Task
		if err := global.DB.WithContext(ctx).Model(&model.Task{}).
			Where("task_type = ? AND deleted_at = ? AND other_id IN ?", service.TASK_TYPE_ARTICLE, 0, aids).
			Order("id ASC").
			Find(&tasks).Error; err != nil {
			return nil, err
		}
		for _, t := range tasks {
			articleByAid[t.OtherId] = t
		}
	}

	courseIDs := make(map[string]struct{})
	for _, t := range titleArticles {
		courseIDs[t.TaskPid] = struct{}{}
	}
	for _, t := range contentArticles {
		courseIDs[t.TaskPid] = struct{}{}
	}
	for _, t := range articleByAid {
		courseIDs[t.TaskPid] = struct{}{}
	}
	missing := make([]string, 0, len(courseIDs))
	for id := range courseIDs {
		if _, ok := courseMap[id]; !ok && id != "" {
			missing = append(missing, id)
		}
	}
	if len(missing) > 0 {
		var courses []*model.Task
		if err := global.DB.WithContext(ctx).Model(&model.Task{}).
			Where("task_type = ? AND deleted_at = ? AND task_id IN ?", service.TASK_TYPE_PRODUCT, 0, missing).
			Find(&courses).Error; err != nil {
			return nil, err
		}
		for _, t := range courses {
			courseMap[t.TaskId] = t
		}
	}

	agg := newSearchAggregator(keyword)
	for _, t := range matchedCourses {
		agg.addCourse(t, lowerKeyword)
	}
	articleHits := 0
	for _, t := range titleArticles {
		if agg.addArticle(t, lowerKeyword, false) {
			articleHits++
		}
	}
	for _, t := range contentArticles {
		if agg.addArticle(t, lowerKeyword, true) {
			articleHits++
		}
	}
	commentHits := 0
	for _, c := range comments {
		var raw geek.ArticleComment
		if len(c.Raw) > 0 {
			_ = json.Unmarshal(c.Raw, &raw)
		}
		if !strings.Contains(strings.ToLower(raw.CommentContent), lowerKeyword) {
			continue
		}
		article, ok := articleByAid[strconv.FormatInt(c.Aid, 10)]
		if !ok {
			continue
		}
		course, ok := courseMap[article.TaskPid]
		if !ok {
			continue
		}
		titled := firstNonEmpty(article.TaskName, "文章")
		if agg.addMatch(course, SearchMatch{
			Type:      "comment",
			ArticleID: article.TaskId,
			Aid:       article.OtherId,
			Title:     titled,
			Snippet:   snippet(raw.CommentContent, keyword, 80),
		}, 1) {
			commentHits++
		}
	}

	courses := agg.build()
	total := len(courses)
	start := (page - 1) * pageSize
	if start > total {
		start = total
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	return &SearchResult{
		Keyword:          keyword,
		TotalCourses:     total,
		TotalArticleHits: articleHits,
		TotalCommentHits: commentHits,
		Page:             page,
		PageSize:         pageSize,
		Courses:          courses[start:end],
		SearchScope:      "本地已缓存课程标题、课程简介、文章标题、文章正文、评论",
		Hint:             "使用 get_course 查看课程目录，使用 get_article 读取正文，使用 list_comments 阅读评论",
	}, nil
}

type searchAggregator struct {
	keyword string
	order   []string
	courses map[string]*SearchCourse
	scores  map[string]int
}

func newSearchAggregator(keyword string) *searchAggregator {
	return &searchAggregator{
		keyword: keyword,
		order:   make([]string, 0, 16),
		courses: make(map[string]*SearchCourse, 16),
		scores:  make(map[string]int, 16),
	}
}

func (a *searchAggregator) ensure(t *model.Task) *SearchCourse {
	if c, ok := a.courses[t.TaskId]; ok {
		return c
	}
	summary := buildCourseSummary(t)
	c := &SearchCourse{
		CourseID: summary.CourseID,
		Title:    summary.Title,
		Subtitle: summary.Subtitle,
		Author:   summary.Author,
		Matches:  make([]SearchMatch, 0, 4),
	}
	a.courses[t.TaskId] = c
	a.order = append(a.order, t.TaskId)
	return c
}

func (a *searchAggregator) addCourse(t *model.Task, lowerKeyword string) {
	c := a.ensure(t)
	if strings.Contains(strings.ToLower(t.TaskName), lowerKeyword) {
		if a.addMatch(t, SearchMatch{
			Type:    "course_title",
			Title:   t.TaskName,
			Snippet: t.TaskName,
		}, 5) {
			return
		}
	}
	var p geek.ProductBase
	if len(t.Raw) > 0 {
		_ = json.Unmarshal(t.Raw, &p)
	}
	for _, text := range []string{p.Intro, p.Subtitle, p.Share.Content} {
		if text != "" && strings.Contains(strings.ToLower(text), lowerKeyword) {
			a.addMatchByCourse(c, SearchMatch{
				Type:    "course_intro",
				Title:   c.Title,
				Snippet: snippet(text, a.keyword, 80),
			}, 2)
			return
		}
	}
}

func (a *searchAggregator) addArticle(t *model.Task, lowerKeyword string, contentMatch bool) bool {
	c, ok := a.courses[t.TaskPid]
	if !ok {
		return false
	}
	var data geek.ArticleData
	if len(t.Raw) > 0 {
		_ = json.Unmarshal(t.Raw, &data)
	}
	title := firstNonEmpty(t.TaskName, data.Info.Title)
	if !contentMatch && strings.Contains(strings.ToLower(title), lowerKeyword) {
		return a.addMatchByCourse(c, SearchMatch{
			Type:      "article_title",
			ArticleID: t.TaskId,
			Aid:       articleAid(t, data.Info),
			Title:     title,
			Snippet:   title,
		}, 3)
	}
	content := data.Info.Content
	if len(content) == 0 {
		content = data.Info.Cshort
	}
	markdown := content
	if len(content) > 0 {
		if converted, err := service.HTMLConvertMarkdown(content); err == nil {
			markdown = converted
		}
	}
	if data.Info.Summary != "" && strings.Contains(strings.ToLower(data.Info.Summary), lowerKeyword) {
		return a.addMatchByCourse(c, SearchMatch{
			Type:      "article_summary",
			ArticleID: t.TaskId,
			Aid:       articleAid(t, data.Info),
			Title:     title,
			Snippet:   snippet(data.Info.Summary, a.keyword, 80),
		}, 2)
	}
	if markdown == "" {
		return false
	}
	if !strings.Contains(strings.ToLower(markdown), lowerKeyword) {
		return false
	}
	return a.addMatchByCourse(c, SearchMatch{
		Type:      "article_content",
		ArticleID: t.TaskId,
		Aid:       articleAid(t, data.Info),
		Title:     title,
		Snippet:   snippet(markdown, a.keyword, 80),
	}, 1)
}

func (a *searchAggregator) addMatch(t *model.Task, match SearchMatch, score int) bool {
	c, ok := a.courses[t.TaskId]
	if !ok {
		c = a.ensure(t)
	}
	return a.addMatchByCourse(c, match, score)
}

func (a *searchAggregator) addMatchByCourse(c *SearchCourse, match SearchMatch, score int) bool {
	if len(c.Matches) >= 8 {
		return false
	}
	for _, m := range c.Matches {
		if m.Type == match.Type && m.ArticleID == match.ArticleID && m.Snippet == match.Snippet {
			return false
		}
	}
	c.Matches = append(c.Matches, match)
	a.scores[c.CourseID] += score
	return true
}

func (a *searchAggregator) build() []SearchCourse {
	out := make([]SearchCourse, 0, len(a.order))
	for _, id := range a.order {
		c := a.courses[id]
		score := a.scores[id]
		if score <= 0 || len(c.Matches) == 0 {
			continue
		}
		c.Score = score
		out = append(out, *c)
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Score == out[j].Score {
			return out[i].Title < out[j].Title
		}
		return out[i].Score > out[j].Score
	})
	return out
}

func buildCourseSummary(t *model.Task) CourseSummary {
	var p geek.ProductBase
	if len(t.Raw) > 0 {
		_ = json.Unmarshal(t.Raw, &p)
	}
	count := p.Article.Count
	if count == 0 {
		count = p.Article.CountPub
	}
	if count == 0 {
		count = p.Article.CountReq
	}
	if count == 0 && len(t.Statistics) > 0 {
		var stats struct {
			Count int `json:"count"`
		}
		if err := json.Unmarshal(t.Statistics, &stats); err == nil {
			count = stats.Count
		}
	}
	return CourseSummary{
		CourseID:     t.TaskId,
		ProductID:    t.OtherId,
		Title:        firstNonEmpty(t.TaskName, p.Title, p.Share.Title),
		Subtitle:     p.Subtitle,
		Author:       p.Author.Name,
		AuthorIntro:  p.Author.Intro,
		Cover:        firstNonEmpty(t.Cover, p.Cover.Square),
		ArticleCount: count,
		IsVideo:      p.IsVideo,
		IsAudio:      p.IsAudio,
		UpdatedAt:    t.UpdatedAt,
	}
}

func findCourseTask(ctx context.Context, id string) (*model.Task, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("course_id 不能为空")
	}
	var t model.Task
	if err := global.DB.WithContext(ctx).Model(&model.Task{}).
		Where("task_type = ? AND deleted_at = ? AND task_id = ?", service.TASK_TYPE_PRODUCT, 0, id).
		First(&t).Error; err == nil {
		return &t, nil
	}
	if err := global.DB.WithContext(ctx).Model(&model.Task{}).
		Where("task_type = ? AND deleted_at = ? AND other_id = ?", service.TASK_TYPE_PRODUCT, 0, id).
		Order("id DESC").
		First(&t).Error; err != nil {
		return nil, fmt.Errorf("课程不存在或未缓存: %s", id)
	}
	return &t, nil
}

func findArticleTask(ctx context.Context, id string) (*model.Task, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("article_id 不能为空")
	}
	var t model.Task
	if err := global.DB.WithContext(ctx).Model(&model.Task{}).
		Where("task_type = ? AND deleted_at = ? AND task_id = ?", service.TASK_TYPE_ARTICLE, 0, id).
		First(&t).Error; err == nil {
		return &t, nil
	}
	if err := global.DB.WithContext(ctx).Model(&model.Task{}).
		Where("task_type = ? AND deleted_at = ? AND other_id = ?", service.TASK_TYPE_ARTICLE, 0, id).
		Order("id DESC").
		First(&t).Error; err != nil {
		return nil, fmt.Errorf("文章不存在或未缓存: %s", id)
	}
	return &t, nil
}

func articleAid(t *model.Task, info geek.ArticleInfo) string {
	if t.OtherId != "" {
		return t.OtherId
	}
	if info.ID > 0 {
		return strconv.Itoa(info.ID)
	}
	return ""
}

func rawTextExpr() string {
	switch strings.ToLower(global.CONF.DB.Driver) {
	case "postgres":
		return "LOWER(CONVERT_FROM(raw, 'UTF8'))"
	case "mysql":
		return "LOWER(CAST(raw AS CHAR))"
	default:
		return "LOWER(CAST(raw AS TEXT))"
	}
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func snippet(text, keyword string, radius int) string {
	text = strings.TrimSpace(text)
	if text == "" {
		return ""
	}
	if radius <= 0 {
		radius = 80
	}
	runes := []rune(text)
	idx := strings.Index(strings.ToLower(text), strings.ToLower(keyword))
	if idx < 0 {
		if len(runes) > radius*2 {
			return strings.TrimSpace(string(runes[:radius*2])) + "…"
		}
		return text
	}
	pos := utf8.RuneCountInString(text[:idx])
	start := pos - radius
	if start < 0 {
		start = 0
	}
	end := pos + utf8.RuneCountInString(keyword) + radius
	if end > len(runes) {
		end = len(runes)
	}
	prefix, suffix := "", ""
	if start > 0 {
		prefix = "…"
	}
	if end < len(runes) {
		suffix = "…"
	}
	return prefix + strings.TrimSpace(string(runes[start:end])) + suffix
}

func truncateRunes(text string, max int) (string, bool) {
	if max <= 0 {
		return text, false
	}
	runes := []rune(text)
	if len(runes) <= max {
		return text, false
	}
	return string(runes[:max]) + "\n\n…(内容已截断)", true
}

func formatUnix(ts int64) string {
	if ts <= 0 {
		return ""
	}
	return time.Unix(ts, 0).In(time.Local).Format("2006-01-02 15:04:05")
}
