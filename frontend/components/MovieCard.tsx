import { useRouter } from 'next/navigation';

interface MovieCardProps {
  id: string;
  title: string;
  genre: string;
  year: number;
  imageUrl: string;
}

export default function MovieCard({ title, genre, year, imageUrl }: MovieCardProps) {
  const router = useRouter();

  const handleClick = () => {
    router.push(`/movie/?title=${title}`);
  };

  return (
    <div
      onClick={handleClick}
      className="rounded-lg shadow-md overflow-hidden cursor-pointer transition-transform hover:scale-105"
    >
      <div
        className="h-48 bg-cover bg-center"
        style={{ backgroundImage: `url(${imageUrl || ''})` }}
      />
      <div className="bg-sky-400 p-2 h-full rounded shadow-md">
        <h3 className="font-semibold text-white">{title}</h3>
        <p className="text-sky-100">{genre} - {year}</p>
      </div>
    </div>
  );
}