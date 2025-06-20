<!-- src/routes/+page.svelte -->
<script>
    import Navbar from '$lib/components/Navbar.svelte';
    import PostCard from '$lib/components/PostCard.svelte';
    
    let activeTab = 'discover';
    let searchQuery = '';
    
    // Mock data - replace with API calls to your Go backend
    const posts = [
      {
        id: 1,
        author: 'The Tennessee Holler',
        handle: '@theholler.bsky.social',
        time: '1d',
        content: 'Another SPACE X explosion',
        image: '/api/placeholder/600/300',
        replies: 1100,
        reposts: 2600,
        likes: 10100,
        hasVideo: true
      },
      {
        id: 2,
        author: 'Megan (she/her)',
        handle: '@americanmegalo.bsky.social',
        time: '1d',
        content: 'With everything LA has gone through in the last few months (and couple of years), let\'s celebrate the hatching of:\n\nTEN CALIFORNIA CONDORS!!!!!!!\n\nCongrats to the team at the Los Angeles Zoo!!!',
        image: '/api/placeholder/400/400',
        replies: 89,
        reposts: 234,
        likes: 1200
      }
    ];
  
    const trendingTopics = [
      'NBA Finals',
      'National Guard', 
      'Beyoncé and Miley',
      'Superman',
      'Juneteenth',
      'SCOTUS'
    ];
  
    // Tab functions
    function setActiveTab(tab) {
      activeTab = tab;
    }
  
    // Navbar event handlers
    function handleNewPost() {
      console.log('New post clicked');
      // TODO: Open compose modal or navigate to compose page
    }
  
    function handleNavigation(event) {
      console.log('Navigation:', event.detail);
      // TODO: Handle navigation with your router
    }
  
    // PostCard event handlers
    function handleReply(event) {
      console.log('Reply to post:', event.detail.postId);
      // TODO: Open reply modal or navigate to reply page
    }
  
    function handleRepost(event) {
      console.log('Repost:', event.detail.postId);
      // TODO: Handle repost logic
    }
  
    function handleLike(event) {
      console.log('Like post:', event.detail.postId);
      // TODO: Handle like/unlike logic
    }
  
    function handleShare(event) {
      console.log('Share post:', event.detail.postId);
      // TODO: Open share modal
    }
  
    function handleAuthorClick(event) {
      console.log('Author clicked:', event.detail.handle);
      // TODO: Navigate to author profile
    }
  
    function handlePlayVideo(event) {
      console.log('Play video:', event.detail);
      // TODO: Open video player
    }
  
    function handleViewMedia(event) {
      console.log('View media:', event.detail);
      // TODO: Open media viewer
    }
  
    // Trending topic handler
    function handleTrendingClick(topic) {
      console.log('Trending topic clicked:', topic);
      // TODO: Navigate to topic page or search
    }
  </script>
  
  <svelte:head>
    <title>SkyClone - Discover</title>
  </svelte:head>
  
  <div class="container">
    <!-- Left Sidebar Navigation -->
    <Navbar 
      on:newPost={handleNewPost}
      on:navigate={handleNavigation}
    />
  
    <!-- Main Content -->
    <main class="main-content">
      <!-- Header Tabs -->
      <header class="header-tabs">
        <button 
          class="tab-btn" 
          class:active={activeTab === 'discover'}
          on:click={() => setActiveTab('discover')}
        >
          Discover
        </button>
        <button 
          class="tab-btn" 
          class:active={activeTab === 'following'}
          on:click={() => setActiveTab('following')}
        >
          Following
        </button>
        <button 
          class="tab-btn" 
          class:active={activeTab === 'video'}
          on:click={() => setActiveTab('video')}
        >
          Video
        </button>
      </header>
  
      <!-- Posts Feed -->
      <div class="feed">
        {#each posts as post (post.id)}
          <PostCard 
            {post}
            on:reply={handleReply}
            on:repost={handleRepost}
            on:like={handleLike}
            on:share={handleShare}
            on:authorClick={handleAuthorClick}
            on:playVideo={handlePlayVideo}
            on:viewMedia={handleViewMedia}
          />
        {/each}
      </div>
    </main>
  
    <!-- Right Sidebar -->
    <aside class="right-sidebar">
      <div class="search-box">
        <input 
          type="text" 
          placeholder="Search" 
          bind:value={searchQuery}
          class="search-input"
        />
      </div>
      
      <div class="sidebar-section">
        <div class="section-tabs">
          <button class="section-tab active">Discover</button>
          <button class="section-tab">Following</button>
          <button class="section-tab">Video</button>
          <button class="section-tab">More feeds</button>
        </div>
      </div>
      
      <div class="sidebar-section">
        <h3 class="section-title">📈 Trending</h3>
        <div class="trending-list">
          {#each trendingTopics as topic}
            <button 
              class="trending-item"
              on:click={() => handleTrendingClick(topic)}
            >
              {topic}
            </button>
          {/each}
        </div>
      </div>
      
      <div class="sidebar-footer">
        <a href="/feedback">Feedback</a>
        <a href="/privacy">Privacy</a>
        <a href="/terms">Terms</a>
        <a href="/help">Help</a>
      </div>
    </aside>
  </div>
  
  <style>
    .container {
      display: grid;
      grid-template-columns: 275px 1fr 350px;
      min-height: 100vh;
      max-width: 1200px;
      margin: 0 auto;
    }
  
    /* Main Content */
    .main-content {
      border-right: 1px solid var(--border-color);
      min-height: 100vh;
      background: var(--bg-primary);
    }
  
    .header-tabs {
      display: flex;
      border-bottom: 1px solid var(--border-color);
      position: sticky;
      top: 0;
      background: var(--bg-primary);
      z-index: 10;
    }
  
    .tab-btn {
      flex: 1;
      padding: 16px;
      background: none;
      color: var(--text-secondary);
      font-size: 15px;
      font-weight: 500;
      border-radius: 0;
      position: relative;
      border: none;
      cursor: pointer;
      transition: all 0.2s ease;
    }
  
    .tab-btn:hover {
      background: var(--bg-secondary);
    }
  
    .tab-btn.active {
      color: var(--text-primary);
      font-weight: 700;
    }
  
    .tab-btn.active::after {
      content: '';
      position: absolute;
      bottom: 0;
      left: 50%;
      transform: translateX(-50%);
      width: 60px;
      height: 4px;
      background: var(--primary-blue);
      border-radius: 2px;
    }
  
    .feed {
      max-width: 100%;
    }
  
    /* Right Sidebar */
    .right-sidebar {
      padding: 20px 16px;
      position: sticky;
      top: 0;
      height: 100vh;
      overflow-y: auto;
      background: var(--bg-primary);
    }
  
    .search-box {
      margin-bottom: 20px;
    }
  
    .search-input {
      width: 100%;
      background: var(--bg-secondary);
      border: 1px solid transparent;
      border-radius: 24px;
      padding: 12px 16px;
      font-size: 15px;
      color: var(--text-primary);
    }
  
    .search-input:focus {
      background: var(--bg-primary);
      border-color: var(--primary-blue);
    }
  
    .search-input::placeholder {
      color: var(--text-secondary);
    }
  
    .sidebar-section {
      background: var(--bg-secondary);
      border-radius: 16px;
      margin-bottom: 20px;
      overflow: hidden;
    }
  
    .section-tabs {
      display: flex;
      flex-wrap: wrap;
      padding: 12px;
      gap: 8px;
    }
  
    .section-tab {
      background: none;
      border: 1px solid var(--border-color);
      color: var(--text-secondary);
      padding: 8px 16px;
      border-radius: 20px;
      font-size: 13px;
      transition: all 0.2s ease;
      cursor: pointer;
    }
  
    .section-tab:hover {
      background: var(--bg-primary);
    }
  
    .section-tab.active {
      background: var(--primary-blue);
      color: white;
      border-color: var(--primary-blue);
    }
  
    .section-title {
      padding: 16px 16px 0 16px;
      margin: 0;
      font-size: 20px;
      font-weight: 800;
      color: var(--text-primary);
    }
  
    .trending-list {
      padding: 4px 16px 16px 16px;
      display: flex;
      flex-direction: column;
      gap: 4px;
    }
  
    .trending-item {
      background: none;
      border: none;
      color: var(--text-primary);
      padding: 8px 0;
      text-align: left;
      font-size: 15px;
      border-radius: 8px;
      transition: background-color 0.2s ease;
      cursor: pointer;
    }
  
    .trending-item:hover {
      background: rgba(0, 0, 0, 0.05);
    }
  
    .sidebar-footer {
      display: flex;
      flex-wrap: wrap;
      gap: 12px;
      margin-top: auto;
      padding-top: 20px;
    }
  
    .sidebar-footer a {
      color: var(--text-secondary);
      text-decoration: none;
      font-size: 13px;
      transition: color 0.2s ease;
    }
  
    .sidebar-footer a:hover {
      text-decoration: underline;
      color: var(--primary-blue);
    }
  
    /* Responsive Design */
    @media (max-width: 1024px) {
      .container {
        grid-template-columns: 68px 1fr 300px;
      }
    }
  
    @media (max-width: 768px) {
      .container {
        grid-template-columns: 1fr;
      }
      
      .right-sidebar {
        display: none;
      }
      
      .main-content {
        border: none;
      }
    }
  </style>