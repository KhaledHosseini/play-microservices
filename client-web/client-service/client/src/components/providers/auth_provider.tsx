'use client';

import {User} from '@/types'
import React from 'react';

const UserContext = React.createContext<
  [User | null, React.Dispatch<React.SetStateAction<User | null>>] | undefined
>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = React.useState<User | null>(null);
  return (
    <UserContext.Provider value={[user, setUser]}>
      {children}
    </UserContext.Provider>
  );
}

export function useAuth() {
  const context = React.useContext(UserContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within a UserContext');
  }
  return context;
}