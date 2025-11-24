package front

import (
	"html/template"
	"net/http"
)

// Handler for serving the web interface
func WebInterfaceHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>Comment Tree</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .comment { margin: 10px 0; padding: 10px; border: 1px solid #ccc; border-radius: 4px; }
        .children { margin-left: 30px; }
        .parent { font-weight: bold; }
        .input-section { margin: 20px 0; }
        .search-box { margin: 20px 0; }
        .search-box input { width: 300px; padding: 5px; }
        .search-box button { padding: 5px 10px; }
        .comment-form { margin-top: 10px; }
        .comment-form textarea { width: 100%; height: 80px; margin: 5px 0; }
        .comment-form button { padding: 5px 10px; }
        .delete-btn { background-color: #ff6b6b; color: white; border: none; padding: 3px 6px; cursor: pointer; }
        .reply-btn { background-color: #4ecdc4; color: white; border: none; padding: 3px 6px; cursor: pointer; }
    </style>
</head>
<body>
    <h1>Comment Tree</h1>

    <div class="search-box">
        <input type="text" id="searchInput" placeholder="Search comments...">
        <button onclick="searchComments()">Search</button>
        <button onclick="loadComments()">Reset</button>
    </div>

    <div class="input-section">
        <h3>Add Root Comment</h3>
        <textarea id="newCommentText" placeholder="Enter your comment"></textarea>
        <br>
        <button onclick="addComment(null)">Add Comment</button>
    </div>

    <div id="commentsContainer"></div>

    <script>
        let commentsData = [];

        // Load comments on page load
        window.onload = function() {
            loadComments();
        };

        // Load all comments
        function loadComments() {
            fetch('/comments')
                .then(response => response.json())
                .then(data => {
                    commentsData = data;
                    renderComments(commentsData);
                })
                .catch(error => console.error('Error:', error));
        }

        // Search comments
        function searchComments() {
            const query = document.getElementById('searchInput').value;
            if (!query) {
                loadComments();
                return;
            }

            fetch('/comments/search?q=' + encodeURIComponent(query))
                .then(response => response.json())
                .then(data => {
                    commentsData = data;
                    renderComments(commentsData);
                })
                .catch(error => console.error('Error:', error));
        }

        // Render comments recursively
        function renderComments(comments, containerId = 'commentsContainer') {
            const container = document.getElementById(containerId);
            container.innerHTML = '';

            comments.forEach(comment => {
                const commentDiv = createCommentElement(comment);
                document.getElementById(containerId).appendChild(commentDiv);
            });
        }

        // Create a comment element
        function createCommentElement(comment) {
            const div = document.createElement('div');
            div.className = 'comment';
            div.id = 'comment-' + comment.id;

            // Format date
            const date = new Date(comment.created_at).toLocaleString();

            div.innerHTML = '<div class="comment-content">' +
                '<div><strong>ID: ' + comment.id + '</strong> | ' + date + '</div>' +
                '<div>' + comment.text + '</div>' +
                '</div>' +
                '<div class="actions">' +
                '<button class="reply-btn" onclick="showReplyForm(' + comment.id + ')">Reply</button>' +
                '<button class="delete-btn" onclick="deleteComment(' + comment.id + ')">Delete</button>' +
                '</div>' +
                '<div class="reply-form" id="reply-form-' + comment.id + '" style="display:none;">' +
                '<textarea placeholder="Enter your reply"></textarea>' +
                '<button onclick="addComment(' + comment.id + ')">Add Reply</button>' +
                '<button onclick="hideReplyForm(' + comment.id + ')">Cancel</button>' +
                '</div>' +
                '<div class="children" id="children-' + comment.id + '"></div>';

            // Render children if any
            if (comment.children && comment.children.length > 0) {
                setTimeout(() => {
                    renderComments(comment.children, 'children-' + comment.id);
                }, 0);
            }

            return div;
        }

        // Show reply form
        function showReplyForm(commentId) {
            document.getElementById('reply-form-' + commentId).style.display = 'block';
        }

        // Hide reply form
        function hideReplyForm(commentId) {
            const form = document.getElementById('reply-form-' + commentId);
            form.style.display = 'none';
            form.querySelector('textarea').value = '';
        }

        // Add a comment
        function addComment(parentId) {
            let text;
            if (parentId !== null) {
                // Reply to a comment
                const form = document.getElementById('reply-form-' + parentId);
                text = form.querySelector('textarea').value;
            } else {
                // Root comment
                text = document.getElementById('newCommentText').value;
            }

            if (!text.trim()) {
                alert('Comment text is required');
                return;
            }

            fetch('/comments', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    parent_id: parentId,
                    text: text
                })
            })
            .then(response => response.json())
            .then(data => {
                if (parentId !== null) {
                    hideReplyForm(parentId);
                } else {
                    document.getElementById('newCommentText').value = '';
                }
                loadComments(); // Refresh the list
            })
            .catch(error => console.error('Error:', error));
        }

        // Delete a comment
        function deleteComment(commentId) {
            if (!confirm('Are you sure you want to delete this comment and all its replies?')) {
                return;
            }

            fetch('/comments/' + commentId, {
                method: 'DELETE'
            })
            .then(response => {
                if (response.ok) {
                    loadComments(); // Refresh the list
                } else {
                    alert('Error deleting comment');
                }
            })
            .catch(error => console.error('Error:', error));
        }
    </script>
</body>
</html>
`

	t, err := template.New("web").Parse(tmpl)
	if err != nil {
		http.Error(w, "Template parsing error", http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, nil)
	if err != nil {
		http.Error(w, "Template execution error", http.StatusInternalServerError)
	}
}
