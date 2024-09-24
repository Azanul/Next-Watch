"use client"

import { useState, useEffect } from "react";
import { useRouter } from 'next/navigation';
import GoogleButton from "@/components/googleButton";
import MainComponent from "@/components/MainComponent";

export default function Home() {
  const router = useRouter();
  const [hasAccessToken, setHasAccessToken] = useState(false);

  useEffect(() => {
    const cookies = document.cookie.split(';');
    const hasToken = cookies.some(cookie => cookie.trim().startsWith('access_token='));
    setHasAccessToken(hasToken);
  }, []);

  function handleGoogleSignIn() {
    router.push(`http://localhost:8080/auth/signin/google`);
  }

  return (
    <div className="grid grid-rows-[20px_1fr_20px] items-center justify-items-center min-h-screen p-8 pb-20 gap-16 sm:p-20 font-[family-name:var(--font-geist-sans)] bg-gradient-to-b from-sky-400 to-sky-200">
      <main className="flex gap-8 row-start-2 w-full h-full items-center place-content-center">
        {hasAccessToken ? (
          <MainComponent />
        ) : (
          <div className="bg-gray-200 w-1/3 h-1/2 place-content-center rounded-md">
            <GoogleButton onClick={handleGoogleSignIn} />
          </div>
        )}
      </main>
      <footer className="row-start-3 flex gap-6 flex-wrap items-center justify-center">
      </footer>
    </div>
  );
}