"use client"

import GoogleButton from "@/components/googleButton";
import { useRouter } from 'next/navigation'

export default function Home() {
  const router = useRouter()

  function handleGoogleSignIn(): void {
    router.push(`http://localhost:8080/auth/signin/google`);
  }

  return (
    <div className="grid grid-rows-[20px_1fr_20px] items-center justify-items-center min-h-screen p-8 pb-20 gap-16 sm:p-20 font-[family-name:var(--font-geist-sans)]">
      <main className="flex gap-8 row-start-2 w-full h-full items-center place-content-center">
        <div className="bg-gray-200 w-1/3 h-1/2 place-content-center rounded-md">
          <GoogleButton onClick={handleGoogleSignIn} />
        </div>
      </main>
      <footer className="row-start-3 flex gap-6 flex-wrap items-center justify-center">
      </footer>
    </div>
  );
}
