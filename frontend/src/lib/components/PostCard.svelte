<!-- PostCard.svelte -->
<script>
  import { deletePost } from '$lib/services/api';

  export let post = {};
  export let user = null;

  let liked = post.user_liked || false;
  let likeCount = post.likes_count || 0;
  let showComments = false;
  let newComment = '';
  let comments = post.comments || [];
  let errorMessage = '';

  function toggleLike() {
    // Placeholder: Backend endpoint needed for likes
    liked = !liked;
    likeCount += liked ? 1 : -1;
    console.log('Toggling like for post:', post.id);
  }

  function toggleComments() {
    showComments = !showComments;
  }

  function addComment() {
    if (!newComment.trim()) return;
    
    comments = [...comments, {
      id: Date.now(),
      author: user?.name || 'Current User',
      text: newComment.trim(),
      timestamp: 'now'
    }];
    newComment = '';
    console.log('Adding comment to post:', post.id);
  }

  async function deletePostHandler() {
    if (!confirm('Are you sure you want to delete this post?')) return;
    
    try {
      await deletePost(post.id);
    } catch (error) {
      errorMessage = 'Failed to delete post. Please try again.';
      console.error('Delete post error:', error);
    }
  }

  function sharePost() {
    console.log('Sharing post:', post.id);
  }

  function getInitials(name) {
    return name?.split(' ').map(n => n.charAt(0)).join('').substring(0, 2) || 'U';
  }

  function formatTime(timestamp) {
    if (timestamp === 'now') return 'now';
    if (!timestamp) return '10m';
    const date = new Date(timestamp * 1000); // Unix timestamp in seconds
    const now = new Date();
    const diff = (now - date) / 1000; // Seconds
    if (diff < 60) return `${Math.floor(diff)}s`;
    if (diff < 3600) return `${Math.floor(diff / 60)}m`;
    if (diff < 86400) return `${Math.floor(diff / 3600)}h`;
    return date.toLocaleString();
  }
</script>

<article class="post-card">
  <header class="post-header">
    <div class="post-author">
      <div class="author-avatar">
        {getInitials(user?.name && post.user_id === user.id ? user.name : 'User')}
      </div>
      <div class="author-info">
        <h4 class="author-name">{user?.name && post.user_id === user.id ? user.name : 'User'}</h4>
        <div class="post-meta">
          <span class="post-time">{formatTime(post.created_at)}</span>
          <span class="visibility-icon">{post.privacy === 'public' ? '🌐' : post.privacy === 'private' ? '🔒' : '👥'}</span>
        </div>
      </div>
    </div>
    {#if post.user_id === user?.id}
      <button class="post-options" on:click={deletePostHandler}>
        🗑️
      </button>
    {/if}
  </header>

  <div class="post-content">
    {#if post.content}
      <p class="post-text">{post.content}</p>
    {/if}
  </div>

  {#if likeCount > 0 || comments.length > 0}
    <div class="post-stats">
      {#if likeCount > 0}
        <div class="likes-count">
          <span class="like-icon">👍</span>
          <span>{likeCount}</span>
        </div>
      {/if}
      {#if comments.length > 0}
        <button class="comments-count" on:click={toggleComments}>
          {comments.length} comment{comments.length !== 1 ? 's' : ''}
        </button>
      {/if}
    </div>
  {/if}

  <div class="post-actions">
    <button class="action-btn {liked ? 'liked' : ''}" on:click={toggleLike}>
      <span class="action-icon">👍</span>
      <span>Like</span>
    </button>
    <button class="action-btn" on:click={toggleComments}>
      <span class="action-icon">💬</span>
      <span>Comment</span>
    </button>
    <button class="action-btn" on:click={sharePost}>
      <span class="action-icon">📤</span>
      <span>Share</span>
    </button>
  </div>

  {#if showComments}
    <div class="comments-section">
      {#if errorMessage}
        <div class="error-message">{errorMessage}</div>
      {/if}
      <div class="add-comment">
        <div class="comment-avatar">
          {getInitials(user?.name || 'U')}
        </div>
        <div class="comment-input-container">
          <input 
            type="text" 
            placeholder="Write a comment..." 
            bind:value={newComment}
            on:keydown={(e) => e.key === 'Enter' && addComment()}
            class="comment-input"
          />
          <button 
            class="comment-submit" 
            on:click={addComment}
            disabled={!newComment.trim()}
          >
            📤
          </button>
        </div>
      </div>
      {#each comments as comment}
        <div class="comment">
          <div class="comment-avatar">
            {getInitials(comment.author)}
          </div>
          <div class="comment-content">
            <div class="comment-bubble">
              <strong class="comment-author">{comment.author}</strong>
              <p class="comment-text">{comment.text}</p>
            </div>
            <div class="comment-actions">
              <button class="comment-action">Like</button>
              <button class="comment-action">Reply</button>
              <span class="comment-time">{formatTime(comment.timestamp)}</span>
            </div>
          </div>
        </div>
      {/each}
    </div>
  {/if}
</article>

<style>
  .post-card {
    background: white;
    border-radius: 8px;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
    margin-bottom: 16px;
    overflow: hidden;
  }
  
  .post-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 16px;
  }
  
  .post-author {
    display: flex;
    align-items: center;
  }
  .error-message {
    color: red;
    font-size: 0.9em;
    margin-bottom: 1em;
  }
  .author-avatar {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    background: #1877f2;
    color: white;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: bold;
    margin-right: 12px;
  }
  
  .author-name {
    font-size: 15px;
    font-weight: 600;
    color: #1c1e21;
    margin: 0 0 2px 0;
  }
  
  .post-meta {
    display: flex;
    align-items: center;
    gap: 4px;
  }
  
  .post-time {
    font-size: 13px;
    color: #65676b;
  }
  
  .visibility-icon {
    font-size: 12px;
    color: #65676b;
  }
  
  .post-options {
    width: 36px;
    height: 36px;
    border-radius: 50%;
    background: none;
    border: none;
    color: #65676b;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: background-color 0.2s;
  }
  
  .post-options:hover {
    background: #f0f2f5;
  }
  
  .post-content {
    padding: 0 16px;
  }
  
  .post-text {
    font-size: 15px;
    line-height: 1.33;
    color: #1c1e21;
    margin: 0 0 12px 0;
    word-wrap: break-word;
  }
  
  .post-stats {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 8px 16px;
    border-bottom: 1px solid #dadde1;
  }
  
  .likes-count {
    display: flex;
    align-items: center;
    gap: 4px;
    font-size: 13px;
    color: #65676b;
  }
  
  .like-icon {
    width: 18px;
    height: 18px;
    background: #1877f2;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 10px;
  }
  
  .comments-count {
    background: none;
    border: none;
    font-size: 13px;
    color: #65676b;
    cursor: pointer;
    text-decoration: underline;
  }
  
  .post-actions {
    display: flex;
    padding: 4px 0;
    border-top: 1px solid #dadde1;
  }
  
  .action-btn {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    padding: 8px;
    background: none;
    border: none;
    color: #65676b;
    font-size: 15px;
    font-weight: 600;
    cursor: pointer;
    border-radius: 4px;
    margin: 4px;
    transition: background-color 0.2s;
  }
  
  .action-btn:hover {
    background: #f0f2f5;
  }
  
  .action-btn.liked {
    color: #1877f2;
  }
  
  .action-icon {
    font-size: 16px;
  }
  
  .comments-section {
    border-top: 1px solid #dadde1;
    padding: 8px 16px;
  }
  
  .add-comment {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 8px;
  }
  
  .comment-avatar {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    background: #42b883;
    color: white;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 12px;
    font-weight: bold;
    flex-shrink: 0;
  }
  
  .comment-input-container {
    display: flex;
    align-items: center;
    flex: 1;
    background: #f0f2f5;
    border-radius: 20px;
    padding: 8px 12px;
  }
  
  .comment-input {
    flex: 1;
    border: none;
    background: none;
    outline: none;
    font-size: 14px;
    color: #1c1e21;
  }
  
  .comment-input::placeholder {
    color: #65676b;
  }
  
  .comment-submit {
    background: none;
    border: none;
    cursor: pointer;
    font-size: 14px;
    color: #1877f2;
    padding: 0 4px;
    opacity: 0.7;
    transition: opacity 0.2s;
  }
  
  .comment-submit:not(:disabled):hover {
    opacity: 1;
  }
  
  .comment-submit:disabled {
    opacity: 0.3;
    cursor: not-allowed;
  }
  
  .comment {
    display: flex;
    gap: 8px;
    margin-bottom: 8px;
  }
  
  .comment-content {
    flex: 1;
  }
  
  .comment-bubble {
    background: #f0f2f5;
    border-radius: 16px;
    padding: 8px 12px;
    display: inline-block;
    max-width: 100%;
  }
  
  .comment-author {
    font-size: 13px;
    color: #1c1e21;
    display: block;
    margin-bottom: 2px;
  }
  
  .comment-text {
    font-size: 14px;
    color: #1c1e21;
    margin: 0;
    line-height: 1.33;
    word-wrap: break-word;
  }
  
  .comment-actions {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-top: 4px;
    padding-left: 12px;
  }
  
  .comment-action {
    background: none;
    border: none;
    font-size: 12px;
    font-weight: 600;
    color: #65676b;
    cursor: pointer;
    padding: 0;
  }
  
  .comment-action:hover {
    text-decoration: underline;
  }
  
  .comment-time {
    font-size: 12px;
    color: #65676b;
  }
  </style>