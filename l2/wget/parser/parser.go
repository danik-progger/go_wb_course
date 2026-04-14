package parser

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

type LinkRewriter struct {
	baseURL   *url.URL
	outputDir string
	links     []string
	rewritten bool
}

func NewLinkRewriter(baseURL *url.URL, outputDir string) *LinkRewriter {
	return &LinkRewriter{
		baseURL:   baseURL,
		outputDir: outputDir,
	}
}

func (lr *LinkRewriter) ExtractAndRewriteLinks(body io.Reader) ([]string, *bytes.Buffer, error) {
	doc, err := html.Parse(body)
	if err != nil {
		return nil, nil, fmt.Errorf("parsing HTML: %w", err)
	}

	var links []string
	lr.rewriteNode(doc)

	var collector nodeCollector
	collector.collect(doc, lr.baseURL)
	links = collector.links

	var buf bytes.Buffer
	if err := html.Render(&buf, doc); err != nil {
		return nil, nil, fmt.Errorf("rendering HTML: %w", err)
	}

	return links, &buf, nil
}

func (lr *LinkRewriter) rewriteNode(n *html.Node) {
	if n.Type == html.ElementNode {
		switch n.Data {
		case "a", "link":
			for i := range n.Attr {
				if n.Attr[i].Key == "href" {
					resolvedURL, err := lr.resolveURL(n.Attr[i].Val)
					if err != nil {
						continue
					}

					localPath := lr.toLocalPath(resolvedURL)
					if localPath != "" {
						n.Attr[i].Val = localPath
					}
				}
			}
		case "img", "script":
			for i := range n.Attr {
				if n.Attr[i].Key == "src" {
					resolvedURL, err := lr.resolveURL(n.Attr[i].Val)
					if err != nil {
						continue
					}

					localPath := lr.toLocalPath(resolvedURL)
					if localPath != "" {
						n.Attr[i].Val = localPath
					}
				}
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		lr.rewriteNode(c)
	}
}

func (lr *LinkRewriter) resolveURL(href string) (string, error) {
	rel, err := url.Parse(href)
	if err != nil {
		return "", err
	}
	return lr.baseURL.ResolveReference(rel).String(), nil
}

func (lr *LinkRewriter) toLocalPath(resolvedURL string) string {
	u, err := url.Parse(resolvedURL)
	if err != nil {
		return ""
	}

	if u.Host != lr.baseURL.Host {
		return ""
	}

	localPath := u.Path
	if strings.HasSuffix(localPath, "/") || localPath == "" {
		localPath = localPath + "index.html"
	}

	localPath = strings.TrimPrefix(localPath, "/")

	return localPath
}

type nodeCollector struct {
	links []string
}

func (nc *nodeCollector) collect(n *html.Node, baseURL *url.URL) {
	if n.Type == html.ElementNode {
		switch n.Data {
		case "a", "link":
			for _, a := range n.Attr {
				if a.Key == "href" {
					resolvedURL, err := resolveURL(baseURL, a.Val)
					if err != nil {
						continue
					}
					nc.links = append(nc.links, resolvedURL)
				}
			}
		case "img", "script":
			for _, a := range n.Attr {
				if a.Key == "src" {
					resolvedURL, err := resolveURL(baseURL, a.Val)
					if err != nil {
						continue
					}
					nc.links = append(nc.links, resolvedURL)
				}
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		nc.collect(c, baseURL)
	}
}

func resolveURL(base *url.URL, href string) (string, error) {
	rel, err := url.Parse(href)
	if err != nil {
		return "", err
	}
	return base.ResolveReference(rel).String(), nil
}
