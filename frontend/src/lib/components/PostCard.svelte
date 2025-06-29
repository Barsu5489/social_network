<script>
  import { likePost, unlikePost, addComment, getPostComments } from '../services/postService';
  import { user } from '../services/auth/authService';
  import { goto } from '$app/navigation';

  export let post;

  let isLiked = false;
  let likesCount = post.likes_count || 0;
  let showComments = false;
  let newCommentText = '';
  let comments = [];
  let commentsCount = post.comments_count || 0;

  // Check if the current user has liked the post
  $: if ($user && post.liked_by_users) {
    isLiked = post.liked_by_users.includes($user.id);
  }

  async function handleLike(event) {
    event.stopPropagation(); // Prevent card click when liking
    const likeButton = event.currentTarget;

    // Add animation
    likeButton.classList.add('animate-like');
    setTimeout(() => {
      likeButton.classList.remove('animate-like');
    }, 600);

    try {
      if (isLiked) {
        await unlikePost(post.id);
        likesCount--;
      } else {
        await likePost(post.id);
        likesCount++;
      }
      isLiked = !isLiked;
    } catch (error) {
      console.error('Error liking/unliking post:', error);
    }
  }

  async function handleAddComment(event) {
    event.stopPropagation();
    if (newCommentText.trim() === '') return;
    try {
      await addComment(post.id, { content: newCommentText });
      newCommentText = '';
      commentsCount++;
      await fetchComments();
    } catch (error) {
      console.error('Error adding comment:', error);
    }
  }

  async function fetchComments() {
    try {
      comments = await getPostComments(post.id);
    } catch (error) {
      console.error('Error fetching comments:', error);
    }
  }

  function toggleComments(event) {
    event.stopPropagation();
    showComments = !showComments;
    if (showComments && comments.length === 0) {
      fetchComments();
    }
  }

  function navigateToPost() {
    goto(`/post/${post.id}`);
  }
</script>

<div class="post-card" on:click={navigateToPost}>
  <div class="post-header">
    <img src={post.author_avatar_url || "/placeholder-avatar.png"} alt="User Avatar" class="avatar" />
    <div class="post-info">
      <span class="username">{post.author_nickname || 'Anonymous'}</span>
      <span class="timestamp">{new Date(post.created_at * 1000).toLocaleString()}</span>
    </div>
  </div>
  <div class="post-content">
    <p>{post.content}</p>
  </div>
  <div class="post-actions">
    <button on:click={handleLike} class="action-button {isLiked ? 'liked' : ''}">
      <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path d="M7 10v12"></path>
        <path d="M15 5.88 14 10h5.83a2 2 0 0 1 1.92 2.56l-2.33 8A2 2 0 0 1 17.5 22H4a2 2 0 0 1-2-2v-8a2 2 0 0 1 2-2h3z"></path>
      </svg>
      <span>{likesCount}</span>
    </button>
    <button on:click={toggleComments} class="action-button">
      <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
      </svg>
      <span>{commentsCount}</span>
    </button>
  </div>

  {#if showComments}
    <div class="comments-section" on:click|stopPropagation>
      <div class="comment-input">
        <input
          type="text"
          placeholder="Add a comment..."
          bind:value={newCommentText}
          on:keydown={(e) => {
            if (e.key === 'Enter') handleAddComment(e);
          }}
        />
        <button on:click={handleAddComment}>Post</button>
      </div>
      <div class="comments-list">
        {#each comments as comment (comment.id)}
          <div class="comment-item">
            <span class="comment-author">{comment.author_username || 'Anonymous'}:</span>
            <span class="comment-content">{comment.content}</span>
          </div>
        {/each}
      </div>
    </div>
  {/if}
</div>

<style>
  .post-card {
    background-color: var(--card-background);
    border-radius: 12px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
    padding: 24px;
    margin-bottom: 20px;
    border: 1px solid var(--border-color);
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .post-card:hover {
    transform: translateY(-3px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    border-color: var(--primary-purple);
  }

  .post-header {
    display: flex;
    align-items: center;
    margin-bottom: 16px;
  }

  .avatar {
    width: 48px;
    height: 48px;
    border-radius: 50%;
    margin-right: 12px;
    background-color: #ccc; /* Placeholder background */
  }

  .post-info .username {
    font-weight: 600;
    color: var(--text-color-primary);
  }

  .post-info .timestamp {
    font-size: 0.8em;
    color: var(--text-color-secondary);
  }

  .post-content {
    margin-bottom: 16px;
    color: var(--text-color-primary);
    font-size: 1.05em;
  }

  .post-actions {
    display: flex;
    gap: 16px;
    border-top: 1px solid var(--border-color);
    padding-top: 12px;
  }

  .action-button {
    background: none;
    border: none;
    color: var(--text-color-secondary);
    font-size: 0.95em;
    font-weight: 500;
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 8px 12px;
    border-radius: 8px;
    transition: all 0.2s ease;
  }

  .action-button:hover {
    background-color: var(--button-hover-background);
    color: var(--primary-purple);
  }

  .action-button.liked {
    color: var(--primary-purple);
  }

  .action-button.animate-like {
    animation: like-animation 0.6s ease-in-out;
  }

  @keyframes like-animation {
    0% { transform: scale(1); }
    50% { transform: scale(1.3); }
    100% { transform: scale(1); }
  }

  .comments-section {
    margin-top: 16px;
    border-top: 1px solid var(--border-color);
    padding-top: 16px;
  }

  .comment-input {
    display: flex;
    gap: 10px;
    margin-bottom: 12px;
  }

  .comment-input input {
    flex-grow: 1;
    padding: 10px 16px;
    border: 1px solid var(--input-border);
    border-radius: 20px;
    font-size: 0.95em;
    background: var(--input-background);
    color: var(--text-color-primary);
  }

  .comment-input button {
    background-color: var(--primary-purple);
    color: white;
    border: none;
    padding: 10px 20px;
    border-radius: 20px;
    font-weight: 600;
    transition: background-color 0.2s ease;
  }

  .comment-input button:hover {
    background-color: var(--deep-purple);
  }

  .comments-list {
    display: flex;
    flex-direction: column;
    gap: 10px;
  }

  .comment-item {
    background-color: var(--button-hover-background);
    padding: 10px 14px;
    border-radius: 15px;
  }

  .comment-author {
    font-weight: 600;
    margin-right: 6px;
    color: var(--text-color-primary);
  }

  .comment-content {
    color: var(--text-color-primary);
  }
</style>