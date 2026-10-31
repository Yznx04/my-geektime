package mcp

import (
	"context"

	mcpgo "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerTools(s *server.MCPServer) {
	s.AddTool(
		newReadOnlyTool("list_courses", "列出本地已缓存的课程（即“我的课程”），按更新时间倒序分页返回。",
			mcpgo.WithNumber("page",
				mcpgo.Description("页码，从 1 开始，默认 1"),
				mcpgo.DefaultNumber(1),
			),
			mcpgo.WithNumber("page_size",
				mcpgo.Description("每页数量，默认 20，最大 100"),
				mcpgo.DefaultNumber(20),
			),
		),
		handleListCourses,
	)

	s.AddTool(
		newReadOnlyTool("search_courses", "按关键词全文检索本地课程：匹配课程标题、课程简介、文章标题、文章正文和评论，返回命中的课程及内容片段。",
			mcpgo.WithString("keyword",
				mcpgo.Description("搜索关键词"),
				mcpgo.Required(),
			),
			mcpgo.WithNumber("page",
				mcpgo.Description("页码，从 1 开始，默认 1"),
				mcpgo.DefaultNumber(1),
			),
			mcpgo.WithNumber("page_size",
				mcpgo.Description("每页课程数量，默认 20，最大 100"),
				mcpgo.DefaultNumber(20),
			),
		),
		handleSearchCourses,
	)

	s.AddTool(
		newReadOnlyTool("get_course", "获取课程详情和章节目录。course_id 可以是列表返回的 task_id，也可以是极客时间课程 ID。",
			mcpgo.WithString("course_id",
				mcpgo.Description("课程 ID（task_id 或极客时间 product id）"),
				mcpgo.Required(),
			),
		),
		handleGetCourse,
	)

	s.AddTool(
		newReadOnlyTool("get_article", "获取指定文章的完整正文，默认转换为 Markdown。article_id 可以是章节目录中的 article_id（task_id），也可以是极客时间文章 ID。",
			mcpgo.WithString("article_id",
				mcpgo.Description("文章 ID（task_id 或极客时间 article id）"),
				mcpgo.Required(),
			),
			mcpgo.WithNumber("max_chars",
				mcpgo.Description("正文最大字符数，0 表示不限制，默认 0"),
				mcpgo.DefaultNumber(0),
			),
		),
		handleGetArticle,
	)

	s.AddTool(
		newReadOnlyTool("list_comments", "读取指定文章的评论；include_replies=true 时会附带本地缓存的回复讨论。",
			mcpgo.WithString("article_id",
				mcpgo.Description("文章 ID（task_id 或极客时间 article id）"),
				mcpgo.Required(),
			),
			mcpgo.WithNumber("page",
				mcpgo.Description("页码，从 1 开始，默认 1"),
				mcpgo.DefaultNumber(1),
			),
			mcpgo.WithNumber("page_size",
				mcpgo.Description("每页评论数量，默认 20，最大 100"),
				mcpgo.DefaultNumber(20),
			),
			mcpgo.WithBoolean("include_replies",
				mcpgo.Description("是否返回评论的回复，默认 false"),
				mcpgo.DefaultBool(false),
			),
		),
		handleListComments,
	)
}

func handleListCourses(ctx context.Context, req mcpgo.CallToolRequest) (*mcpgo.CallToolResult, error) {
	page := clamp(req.GetInt("page", 1), 1, 1_000_000)
	pageSize := clamp(req.GetInt("page_size", 20), 1, 100)
	ret, err := listCourses(ctx, page, pageSize)
	if err != nil {
		return mcpgo.NewToolResultError(err.Error()), nil
	}
	return mcpgo.NewToolResultJSON(ret)
}

func handleSearchCourses(ctx context.Context, req mcpgo.CallToolRequest) (*mcpgo.CallToolResult, error) {
	keyword, err := req.RequireString("keyword")
	if err != nil {
		return mcpgo.NewToolResultError("缺少参数 keyword"), nil
	}
	page := clamp(req.GetInt("page", 1), 1, 1_000_000)
	pageSize := clamp(req.GetInt("page_size", 20), 1, 100)
	ret, err := searchCourses(ctx, keyword, page, pageSize)
	if err != nil {
		return mcpgo.NewToolResultError(err.Error()), nil
	}
	return mcpgo.NewToolResultJSON(ret)
}

func handleGetCourse(ctx context.Context, req mcpgo.CallToolRequest) (*mcpgo.CallToolResult, error) {
	courseID, err := req.RequireString("course_id")
	if err != nil {
		return mcpgo.NewToolResultError("缺少参数 course_id"), nil
	}
	ret, err := getCourse(ctx, courseID)
	if err != nil {
		return mcpgo.NewToolResultError(err.Error()), nil
	}
	return mcpgo.NewToolResultJSON(ret)
}

func handleGetArticle(ctx context.Context, req mcpgo.CallToolRequest) (*mcpgo.CallToolResult, error) {
	articleID, err := req.RequireString("article_id")
	if err != nil {
		return mcpgo.NewToolResultError("缺少参数 article_id"), nil
	}
	maxChars := req.GetInt("max_chars", 0)
	if maxChars < 0 {
		maxChars = 0
	}
	ret, err := getArticle(ctx, articleID, maxChars)
	if err != nil {
		return mcpgo.NewToolResultError(err.Error()), nil
	}
	return mcpgo.NewToolResultJSON(ret)
}

func handleListComments(ctx context.Context, req mcpgo.CallToolRequest) (*mcpgo.CallToolResult, error) {
	articleID, err := req.RequireString("article_id")
	if err != nil {
		return mcpgo.NewToolResultError("缺少参数 article_id"), nil
	}
	page := clamp(req.GetInt("page", 1), 1, 1_000_000)
	pageSize := clamp(req.GetInt("page_size", 20), 1, 100)
	ret, err := listComments(ctx, articleID, page, pageSize, req.GetBool("include_replies", false))
	if err != nil {
		return mcpgo.NewToolResultError(err.Error()), nil
	}
	return mcpgo.NewToolResultJSON(ret)
}

func newReadOnlyTool(name, description string, opts ...mcpgo.ToolOption) mcpgo.Tool {
	base := []mcpgo.ToolOption{
		mcpgo.WithDescription(description),
		mcpgo.WithReadOnlyHintAnnotation(true),
		mcpgo.WithDestructiveHintAnnotation(false),
		mcpgo.WithIdempotentHintAnnotation(true),
		mcpgo.WithOpenWorldHintAnnotation(false),
	}
	return mcpgo.NewTool(name, append(base, opts...)...)
}

func clamp(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
