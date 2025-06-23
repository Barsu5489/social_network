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
