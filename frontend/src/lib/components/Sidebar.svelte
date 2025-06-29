<script>
  // This is a placeholder for the sidebar content
</script>

<div class="sidebar-container">
  <h3>Categories</h3>
  <ul>
    <li><a href="#">Friends</a></li>
    <li><a href="#">Groups</a></li>
    <li><a href="#">Events</a></li>
    <li><a href="#">Marketplace</a></li>
  </ul>

  <h3>Your Shortcuts</h3>
  <ul>
    <li><a href="#">Photos</a></li>
    <li><a href="#">Videos</a></li>
  </ul>
</div>

<style>
  .sidebar-container {
    background-color: var(--background-white);
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    padding: 20px;
    height: fit-content;
  }

  .sidebar-container h3 {
    color: var(--text-primary);
    font-size: 1.1em;
    margin-top: 0;
    margin-bottom: 15px;
  }

  .sidebar-container ul {
    list-style: none;
    padding: 0;
    margin-bottom: 20px;
  }

  .sidebar-container ul li {
    margin-bottom: 10px;
  }

  .sidebar-container ul li a {
    color: var(--text-primary);
    text-decoration: none;
    display: block;
    padding: 5px 0;
  }

  .sidebar-container ul li a:hover {
    color: var(--primary-purple);
    text-decoration: underline;
  }
</style>
