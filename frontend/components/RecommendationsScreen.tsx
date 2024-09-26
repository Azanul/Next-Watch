import MoviesComponent from "./MoviesComponent";

export default function RecommendationsScreen() {
  return (
    <div className="container mx-auto px-4">
      <MoviesComponent queryType='GET_RECOMMENDATIONS'></MoviesComponent>
    </div>
  )
}