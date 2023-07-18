import { useAuth } from '../providers/auth_provider'

export default function Header() {
  
  const [user,_setUser] = useAuth()

  return (
    <header className="w-full text-gray-700 bg-white border-t border-gray-100 shadow-sm body-font">
    <div className="container flex flex-col flex-wrap items-center p-5 mx-auto md:flex-row">

        <nav className="flex flex-wrap items-center text-base lg:w-2/5 md:ml-auto">
            <h1 className="mr-5 font-medium hover:text-gray-900">Microservice Job Scheduler</h1>
        </nav>
        
        <div className="inline-flex items-center h-full ml-5 lg:w-2/5 lg:justify-end lg:ml-0">
            {user == null && <a href="/login" className="mr-5 font-medium hover:text-gray-900">Login</a>}
            {user == null && <a href="/signup" className="px-4 py-2 text-xs font-bold text-white uppercase transition-all duration-150 bg-teal-500 rounded shadow outline-none active:bg-teal-600 hover:shadow-md focus:outline-none ease"> Sign Up</a>}
            {user != null && <h2 className="mr-5 font-medium hover:text-gray-900">{user.email}</h2>}
            {user != null && <a href="/logout" className="px-4 py-2 text-xs font-bold text-white uppercase transition-all duration-150 bg-teal-500 rounded shadow outline-none active:bg-teal-600 hover:shadow-md focus:outline-none ease">Log out</a>}
        </div>
    </div>
</header>
  )
}
