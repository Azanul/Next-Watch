"use client"

import localFont from "next/font/local";
import "./globals.css";
import { ApolloWrapper } from "@/components/ApolloWrapper";
import FormbricksProvider from "./formbricks";
import { Suspense, useEffect, useState } from "react";
import Footer from "./footer";

const geistSans = localFont({
  src: "./fonts/GeistVF.woff",
  variable: "--font-geist-sans",
  weight: "100 900",
});

const geistMono = localFont({
  src: "./fonts/GeistMonoVF.woff",
  variable: "--font-geist-mono",
  weight: "100 900",
});

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const [hasAccessToken, setHasAccessToken] = useState(false);

  useEffect(() => {
    const cookies = document.cookie.split(';');
    const hasToken = cookies.some(cookie => cookie.trim().startsWith('access_token='));
    setHasAccessToken(hasToken);
  }, []);

  return (
    <html lang="en">
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased font-[family-name:var(--font-geist-sans)] bg-gradient-to-b from-sky-400 to-sky-200`}
      >
        <Suspense>
          <ApolloWrapper>{children}</ApolloWrapper>
          {hasAccessToken && <FormbricksProvider />}
        </Suspense>
      </body>
      <Footer />
    </html>
  );
}
