<!-- src/lib/components/PostCard.svelte -->
<script>
    import { createEventDispatcher } from 'svelte';
    
    export let post;
    
    const dispatch = createEventDispatcher();
    
    function formatNumber(num) {
      if (num >= 1000000) return (num / 1000000).toFixed(1) + 'M';
      if (num >= 1000) return (num / 1000).toFixed(1) + 'K';
      return num.toString();
    }
  
    function handleReply() {
      dispatch('reply', { postId: post.id });
    }
  
    function handleRepost() {
      dispatch('repost', { postId: post.id });
    }
  
    function handleLike() {
      dispatch('like', { postId: post.id });
    }
  
    function handleShare() {
      dispatch('share', { postId: post.id });
    }
  
    function handleAuthorClick() {
      dispatch('authorClick', { handle: post.handle });
    }
  
    function handleMediaClick() {
      if (post.hasVideo) {
        dispatch('playVideo', { postId: post.id, media: post.image });
      } else {
        dispatch('viewMedia', { postId: post.id, media: post.image });
      }
    }
  </script>
  
  <article class="post">
    <div class="post-header">
      <button class="author-avatar" on:click={handleAuthorClick}></button>
      <div class="author-info">
        <button class="author-name" on:click={handleAuthorClick}>
          {post.author}
        </button>
        <button class="author-handle" on:click={handleAuthorClick}>
          {post.handle}
        </button>
        <span class="post-time">· {post.time}</span>
      </div>
    </div>
    
    <div class="post-content">
      <p>{post.content}</p>
      {#if post.image}
        <div class="post-media" on:click={handleMediaClick} role="button" tabindex="0">
          {#if post.hasVideo}
            <div class="video-overlay">
              <div class="play-button">▶</div>
            </div>
          {/if}
          <img src={post.image} alt="Post media" />
        </div>
      {/if}
    </div>
    
    <div class="post-actions">
      <button class="action-btn reply-btn" on:click={handleReply}>
        <span class="action-icon">💬</span>
        <span class="action-count">{formatNumber(post.replies)}</span>
      </button>
      <button class="action-btn repost-btn" on:click={handleRepost}>
        <span class="action-icon">🔄</span>
        <span class="action-count">{formatNumber(post.reposts)}</span>
      </button>
      <button class="action-btn like-btn" on:click={handleLike}>
        <span class="action-icon">❤️</span>
        <span class="action-count">{formatNumber(post.likes)}</span>
      </button>
      <button class="action-btn share-btn" on:click={handleShare}>
        <span class="action-icon">📤</span>
      </button>
    </div>
  </article>
  
  <style>
    .post {
      padding: 16px 20px;
      border-bottom: 1px solid var(--border-color);
      transition: background-color 0.2s ease;
      cursor: pointer;
    }
  
    .post:hover {
      background: rgba(0, 0, 0, 0.02);
    }
  
    .post-header {
      display: flex;
      gap: 12px;
      margin-bottom: 12px;
    }
  
    .author-avatar {
      width: 40px;
      height: 40px;
      border-radius: 50%;
      background: linear-gradient(135deg, #667eea, #764ba2);
      flex-shrink: 0;
      border: none;
      cursor: pointer;
      transition: transform 0.2s ease;
    }
  
    .author-avatar:hover {
      transform: scale(1.05);
    }
  
    .author-info {
      display: flex;
      align-items: center;
      gap: 8px;
      flex-wrap: wrap;
      flex: 1;
    }
  
    .author-name {
      font-weight: 700;
      font-size: 15px;
      margin: 0;
      color: var(--text-primary);
      background: none;
      border: none;
      cursor: pointer;
      text-align: left;
      padding: 0;
      transition: color 0.2s ease;
    }
  
    .author-name:hover {
      text-decoration: underline;
    }
  
    .author-handle {
      color: var(--text-secondary);
      font-size: 15px;
      background: none;
      border: none;
      cursor: pointer;
      text-align: left;
      padding: 0;
      transition: color 0.2s ease;
    }
  
    .author-handle:hover {
      text-decoration: underline;
    }
  
    .post-time {
      color: var(--text-secondary);
      font-size: 15px;
    }
  
    .post-content p {
      margin: 0 0 12px 0;
      font-size: 15px;
      line-height: 1.5;
      white-space: pre-line;
      color: var(--text-primary);
    }
  
    .post-media {
      position: relative;
      border-radius: 16px;
      overflow: hidden;
      margin-top: 12px;
      cursor: pointer;
      transition: transform 0.2s ease;
    }
  
    .post-media:hover {
      transform: scale(0.98);
    }
  
    .post-media img {
      width: 100%;
      height: auto;
      display: block;
    }
  
    .video-overlay {
      position: absolute;
      top: 50%;
      left: 50%;
      transform: translate(-50%, -50%);
      z-index: 2;
      pointer-events: none;
    }
  
    .play-button {
      width: 60px;
      height: 60px;
      border-radius: 50%;
      background: rgba(0, 0, 0, 0.8);
      color: white;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 20px;
      transition: all 0.2s ease;
    }
  
    .post-media:hover .play-button {
      background: rgba(0, 0, 0, 0.9);
      transform: scale(1.1);
    }
  
    .post-actions {
      display: flex;
      gap: 60px;
      margin-top: 12px;
      padding-left: 52px;
    }
  
    .action-btn {
      display: flex;
      align-items: center;
      gap: 8px;
      background: none;
      color: var(--text-secondary);
      font-size: 13px;
      padding: 8px;
      border-radius: 20px;
      transition: all 0.2s ease;
      border: none;
      cursor: pointer;
    }
  
    .action-btn:hover {
      background: rgba(29, 155, 240, 0.1);
      color: var(--primary-blue);
    }
  
    .reply-btn:hover {
      background: rgba(29, 155, 240, 0.1);
      color: #1d9bf0;
    }
  
    .repost-btn:hover {
      background: rgba(0, 186, 124, 0.1);
      color: #00ba7c;
    }
  
    .like-btn:hover {
      background: rgba(249, 24, 128, 0.1);
      color: #f91880;
    }
  
    .share-btn:hover {
      background: rgba(29, 155, 240, 0.1);
      color: #1d9bf0;
    }
  
    .action-icon {
      font-size: 18px;
      transition: transform 0.2s ease;
    }
  
    .action-btn:hover .action-icon {
      transform: scale(1.1);
    }
  
    .action-count {
      font-size: 13px;
      font-weight: 400;
    }
  
    /* Responsive Design */
    @media (max-width: 768px) {
      .post {
        padding: 12px 16px;
      }
  
      .post-actions {
        gap: 40px;
        padding-left: 0;
        justify-content: space-around;
      }
  
      .action-count {
        display: none;
      }
    }
  </style>