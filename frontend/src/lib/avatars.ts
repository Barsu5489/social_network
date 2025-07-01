// Reliable avatar images using Pravatar.cc (verified working service)
export const defaultAvatarUrls = [
  'https://i.pravatar.cc/128?img=1',
  'https://i.pravatar.cc/128?img=5',
  'https://i.pravatar.cc/128?img=8',
  'https://i.pravatar.cc/128?img=12',
  'https://i.pravatar.cc/128?img=16',
  'https://i.pravatar.cc/128?img=20',
  'https://i.pravatar.cc/128?img=25',
  'https://i.pravatar.cc/128?img=32',
];

export const getRandomAvatar = () => {
  return defaultAvatarUrls[Math.floor(Math.random() * defaultAvatarUrls.length)];
};

// Reliable cover images using Picsum Photos (verified working service)
export const defaultGroupCoverUrls = [
  { 
    url: 'https://picsum.photos/1200/400?random=1', 
    hint: 'nature landscape' 
  },
  { 
    url: 'https://picsum.photos/1200/400?random=2', 
    hint: 'urban architecture' 
  },
  { 
    url: 'https://picsum.photos/1200/400?random=3', 
    hint: 'abstract patterns' 
  },
  { 
    url: 'https://picsum.photos/1200/400?random=4', 
    hint: 'scenic views' 
  },
  { 
    url: 'https://picsum.photos/1200/400?random=5', 
    hint: 'modern design' 
  },
  { 
    url: 'https://picsum.photos/1200/400?random=6', 
    hint: 'artistic composition' 
  },
  { 
    url: 'https://picsum.photos/1200/400?random=7', 
    hint: 'creative space' 
  },
  { 
    url: 'https://picsum.photos/1200/400?random=8', 
    hint: 'inspiring imagery' 
  },
];

// This function will select a cover image deterministically based on the group ID.
export const getGroupCover = (groupId: string): { url: string; hint: string } => {
  // Simple hash function to get an index from the groupId
  let hash = 0;
  if (groupId.length === 0) return defaultGroupCoverUrls[0];
  
  for (let i = 0; i < groupId.length; i++) {
    const char = groupId.charCodeAt(i);
    hash = ((hash << 5) - hash) + char;
    hash = hash & hash; // Convert to 32bit integer
  }
  
  const index = Math.abs(hash) % defaultGroupCoverUrls.length;
  return defaultGroupCoverUrls[index];
};

// Optional: Get a random cover image
export const getRandomGroupCover = (): { url: string; hint: string } => {
  return defaultGroupCoverUrls[Math.floor(Math.random() * defaultGroupCoverUrls.length)];
};