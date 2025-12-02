import { useNavigate, useParams } from "react-router-dom"
import { API_URL } from "./App"

export const ConfirmationPage = () => {
  const { token = '' } = useParams()
  const redirect = useNavigate()

  const handleConfirm = async () => {
    const response = await fetch(`${API_URL}/users/activate/${token}`, {
      method: "PUT"
    })

    if (response.ok) {
      redirect("/")
    } else {
      // handle error
      alert("Failed to confirm token")
    }
  }

  return (
    <div className="min-h-screen bg-gray-900 text-white">
      <nav className="bg-gray-800 border-b border-gray-700 px-6 py-4">
        <div className="max-w-6xl mx-auto">
          <h1 className="text-2xl font-bold text-cyan-400">Forseer</h1>
        </div>
      </nav>

      <div className="max-w-md mx-auto px-6 py-20">
        <div className="bg-gray-800 rounded-lg p-8 border border-gray-700">
          <h2 className="text-2xl font-bold mb-6 text-cyan-400">Confirmation</h2>
          <button
            onClick={handleConfirm}
            className="w-full bg-cyan-500 hover:bg-cyan-600 text-white py-2 rounded font-medium"
          >
            Click to confirm
          </button>
        </div>
      </div>
    </div>
  )
}