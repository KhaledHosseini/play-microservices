import './globals.css';
import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import { QueryProvider } from '@/components/providers/query_provider'
import {AuthProvider} from '@/components/providers/auth_provider'
import ToastProvider from '@/components/providers/toast_provider'

const inter = Inter({ subsets: ['latin'] })

export const metadata: Metadata = {
  title: 'Microservices job scheduler',
  description: 'A simple job scheduler app with microservices architecture.',
}

export default async function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {

  return (
    <html lang="en">
      <body className={`${inter.className}  font-inter antialiased bg-gray-300 text-gray-900 tracking-tight`}>
      <ToastProvider />
        <AuthProvider>
          <QueryProvider>
            {children}
          </QueryProvider>
        </AuthProvider>
      </body>
    </html>
  )
}
