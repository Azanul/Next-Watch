"use client"

import GoogleButton from "@/components/googleButton";
import { useRouter } from 'next/navigation'
import { useEffect, useState } from "react";

const SCOPE = "https%3A//www.googleapis.com/auth/drive.metadata.readonly"
const GOOGLE_CLIENT_ID = "914120788754-9b41keif6f5p7qtcide9ifa7bmkkrs2b.apps.googleusercontent.com"
const REDIRECT_URI = "http%3A//localhost:64139/"

export default function Home() {
  const router = useRouter()

  const [hasAccessToken, setHasAccessToken] = useState("");

  useEffect(() => {
      // Get the access_token from the URL query parameters
      const access_token = getAccessTokenFromUrl()

      // Check if access_token exists
      if (access_token) {
        setHasAccessToken(access_token);
        setSecureCookie("access_token", access_token, 3600 * 24 * 7)
      } else {
        console.log('No access token found in URL');
      }
    
  }, []);

  function handleGoogleSignIn(): void {
    router.push(`https://accounts.google.com/o/oauth2/v2/auth?response_type=token&include_granted_scopes=true&scope=${SCOPE}&client_id=${GOOGLE_CLIENT_ID}&redirect_uri=${REDIRECT_URI}`);
  }

  return (
    <div className="grid grid-rows-[20px_1fr_20px] items-center justify-items-center min-h-screen p-8 pb-20 gap-16 sm:p-20 font-[family-name:var(--font-geist-sans)]">
      <main className="flex gap-8 row-start-2 w-full h-full items-center place-content-center">
        <div className="bg-gray-200 w-1/3 h-1/2 place-content-center rounded-md">
          <GoogleButton onClick={handleGoogleSignIn}/>
        </div>
      </main>
      <footer className="row-start-3 flex gap-6 flex-wrap items-center justify-center">
      </footer>
    </div>
  );
}

// This function extracts the access token from the URL fragment
function getAccessTokenFromUrl() {
  const hashParams = new URLSearchParams(window.location.hash.substring(1)); // Skip the #
  const accessToken = hashParams.get('access_token');
  
  if (accessToken) {
      console.log('Access Token:', accessToken);
      return accessToken;
  } else {
      console.log('Access Token not found in URL');
      return null;
  }
}

// Function to set a cookie with attributes
function setSecureCookie(name: any, value: any, expiresInSeconds: number) {
  const d = new Date();
  d.setTime(d.getTime() + (expiresInSeconds * 1000));
  
  let cookieString = `${name}=${value}; expires=${d.toUTCString()}; path=/; SameSite=Strict;`;
  
  // If on HTTPS, add the Secure flag
  if (window.location.protocol === 'https:') {
      cookieString += ' Secure;';
  }

  document.cookie = cookieString;
}