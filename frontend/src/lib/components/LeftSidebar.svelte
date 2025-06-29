<script>
  // This is a placeholder for the left sidebar content
</script>

<div class="sidebar-container">
  <h3>My Profile</h3>
  <ul>
    <li><a href="/profile">View Profile</a></li>
  </ul>

  <h3>Navigation</h3>
  <ul>
    <li><a href="/home">Home</a></li>
    <li><a href="/friends">Friends</a></li>
    <li><a href="/groups">Groups</a></li>
    <li><a href="/events">Events</a></li>
    <li><a href="/messages">Messages</a></li>
  </ul>
</div>

<style>
  .sidebar-container {
    background-color: var(--background-color);
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    padding: 20px;
  }

  .sidebar-container h3 {
    color: var(--text-color-primary);
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
    color: var(--text-color-primary);
    text-decoration: none;
    display: block;
    padding: 5px 0;
  }

  .sidebar-container ul li a:hover {
    color: var(--primary-purple);
    text-decoration: underline;
  }
</style>
