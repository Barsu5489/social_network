<script>
  const birthdays = [
    { name: 'Onyii Richie Senior', id: 1 }
  ];
  
  const contacts = [
    { name: 'Amina Amina', id: 1, online: true },
    { name: 'Leah Helary', id: 2, online: false },
  ];
  
  let searchContacts = '';
  
  $: filteredContacts = contacts.filter(contact => 
    contact.name.toLowerCase().includes(searchContacts.toLowerCase())
  );
  
  function getInitials(name) {
    return name.split(' ').map(n => n.charAt(0)).join('').substring(0, 2);
  }
</script>

<aside class="right-sidebar">
  <div class="sidebar-content">
    <!-- Birthdays Section -->
    <div class="section">
      <h3 class="section-title">Birthdays</h3>
      {#each birthdays as birthday}
        <div class="birthday-item">
          <div class="birthday-icon">🎁</div>
          <span class="birthday-text">
            <strong>{birthday.name}'s</strong> birthday is today.
          </span>
        </div>
      {/each}
    </div>
    
    <div class="section-divider"></div>
    
    <!-- Contacts Section -->
    <div class="section">
      <div class="contacts-header">
        <h3 class="section-title">Contacts</h3>
        <div class="contacts-actions">
          <button class="action-btn" title="Search">
            🔍
          </button>
          <button class="action-btn" title="More options">
            ⋯
          </button>
        </div>
      </div>
      
      <!-- Contact List -->
      <div class="contacts-list">
        {#each filteredContacts as contact}
          <div class="contact-item">
            <div class="contact-avatar">
              <div class="avatar-circle">
                {getInitials(contact.name)}
              </div>
              {#if contact.online}
                <div class="online-indicator"></div>
              {/if}
            </div>
            <span class="contact-name">{contact.name}</span>
          </div>
        {/each}
      </div>
    </div>
  </div>
</aside>

<style>
  .right-sidebar {
    position: fixed;
    right: 0;
    top: 56px;
    width: 280px;
    height: calc(100vh - 56px);
    background: white;
    overflow-y: auto;
    padding: 16px 0;
    border-left: 1px solid #dadde1;
  }
  
  .sidebar-content {
    padding: 0 16px;
  }
  
  .section {
    margin-bottom: 24px;
  }
  
  .section-title {
    font-size: 17px;
    font-weight: 600;
    color: #65676b;
    margin: 0 0 12px 0;
  }
  
  .birthday-item {
    display: flex;
    align-items: center;
    padding: 8px;
    border-radius: 8px;
    transition: background-color 0.2s;
    cursor: pointer;
  }
  
  .birthday-item:hover {
    background: #f0f2f5;
  }
  
  .birthday-icon {
    width: 32px;
    height: 32px;
    background: linear-gradient(45deg, #ff6b6b, #ffd93d);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-right: 12px;
    font-size: 16px;
  }
  
  .birthday-text {
    font-size: 14px;
    color: #1c1e21;
  }
  
  .section-divider {
    height: 1px;
    background: #dadde1;
    margin: 16px 0;
  }
  
  .contacts-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 12px;
  }
  
  .contacts-actions {
    display: flex;
    gap: 4px;
  }
  
  .action-btn {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    background: none;
    border: none;
    color: #65676b;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 14px;
    transition: background-color 0.2s;
  }
  
  .action-btn:hover {
    background: #f0f2f5;
  }
  
  .contacts-list {
    max-height: 400px;
    overflow-y: auto;
  }
  
  .contact-item {
    display: flex;
    align-items: center;
    padding: 8px;
    border-radius: 8px;
    cursor: pointer;
    transition: background-color 0.2s;
    margin-bottom: 2px;
  }
  
  .contact-item:hover {
    background: #f0f2f5;
  }
  
  .contact-avatar {
    position: relative;
    margin-right: 12px;
  }
  
  .avatar-circle {
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
  }
  
  .online-indicator {
    position: absolute;
    bottom: 0;
    right: 0;
    width: 12px;
    height: 12px;
    background: #42b883;
    border: 2px solid white;
    border-radius: 50%;
  }
  
  .contact-name {
    font-size: 14px;
    color: #1c1e21;
    font-weight: 500;
  }
  
  /* Scrollbar styling */
  .right-sidebar::-webkit-scrollbar,
  .contacts-list::-webkit-scrollbar {
    width: 8px;
  }
  
  .right-sidebar::-webkit-scrollbar-track,
  .contacts-list::-webkit-scrollbar-track {
    background: transparent;
  }
  
  .right-sidebar::-webkit-scrollbar-thumb,
  .contacts-list::-webkit-scrollbar-thumb {
    background: #bcc0c4;
    border-radius: 4px;
  }
  
  .right-sidebar::-webkit-scrollbar-thumb:hover,
  .contacts-list::-webkit-scrollbar-thumb:hover {
    background: #8a8d91;
  }
</style>