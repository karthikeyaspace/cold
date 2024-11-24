import React, { useState, useEffect } from "react";
import { api } from "./axios";

interface ExcelData {
  Sno: string;
  Name: string;
  Email: string;
  Company: string;
  ApplyingPosition: string;
  AdditionalInfo: string;
  ReasonForContact: string;
}

interface EmailContent {
  Subject: string;
  HTML: string;
}

interface GeneratedEmails {
  [key: string]: EmailContent;
}

const App: React.FC = () => {
  const [excelData, setExcelData] = useState<ExcelData[]>([]);
  const [generatedEmails, setGeneratedEmails] = useState<GeneratedEmails>({});
  const [loading, setLoading] = useState<boolean>(true);
  const [processingRows, setProcessingRows] = useState<Set<string>>(new Set());
  const [errorMessages, setErrorMessages] = useState<{ [key: string]: string }>(
    {}
  );
  const [successMessages, setSuccessMessages] = useState<{
    [key: string]: string;
  }>({});

  useEffect(() => {
    fetchExcelData();
  }, []);

  const fetchExcelData = async (): Promise<void> => {
    try {
      const res = await api.get("/data");
      if (res.data.success) {
        setExcelData(res.data.data);
      } else throw new Error("Failed to fetch data");
    } catch (error) {
      setErrorMessages((prev) => ({
        ...prev,
        general: "Failed to fetch excel data",
      }));
    } finally {
      setLoading(false);
    }
  };

  const handleGenerateMail = async (row: ExcelData): Promise<void> => {
    setProcessingRows((prev) => new Set(prev).add(row.Sno));
    setErrorMessages((prev) => ({ ...prev, [row.Sno]: "" }));
    setSuccessMessages((prev) => ({ ...prev, [row.Sno]: "" }));

    try {
      const res = await api.post("/generate", JSON.stringify(row));

      if (res.data.success) {
        setGeneratedEmails((prev: GeneratedEmails) => ({
          ...prev,
          [row.Sno]: {
            Subject: res.data.data.Subject,
            HTML: res.data.data.HTML,
          },
        }));
        setSuccessMessages((prev) => ({
          ...prev,
          [row.Sno]: "Email generated successfully",
        }));
      } else {
        throw new Error("Failed to generate email");
      }
    } catch (error) {
      setErrorMessages((prev) => ({
        ...prev,
        [row.Sno]: "Failed to generate email",
      }));
    } finally {
      setProcessingRows((prev) => {
        const next = new Set(prev);
        next.delete(row.Sno);
        return next;
      });
    }
  };

  const handleSendMail = async (row: ExcelData): Promise<void> => {
    const emailContent = generatedEmails[row.Sno];
    if (!emailContent) {
      setErrorMessages((prev) => ({
        ...prev,
        [row.Sno]: "Please generate email first",
      }));
      return;
    }

    setProcessingRows((prev) => new Set(prev).add(row.Sno));
    setErrorMessages((prev) => ({ ...prev, [row.Sno]: "" }));
    setSuccessMessages((prev) => ({ ...prev, [row.Sno]: "" }));

    try {
      const res = await api.post(
        "/send",
        JSON.stringify({
          To: row.Email,
          Subject: emailContent.Subject,
          HTML: emailContent.HTML,
        })
      );
      if (res.data.success) {
        setSuccessMessages((prev) => ({
          ...prev,
          [row.Sno]: "Email sent successfully",
        }));
      } else {
        throw new Error("Failed to send email");
      }
    } catch (error) {
      setErrorMessages((prev) => ({
        ...prev,
        [row.Sno]: "Failed to send email",
      }));
    } finally {
      setProcessingRows((prev) => {
        const next = new Set(prev);
        next.delete(row.Sno);
        return next;
      });
    }
  };

  const handleEmailContentChange = (
    sno: string,
    field: keyof EmailContent,
    value: string
  ): void => {
    setGeneratedEmails((prev: GeneratedEmails) => ({
      ...prev,
      [sno]: {
        ...prev[sno],
        [field]: value,
      },
    }));
  };

  if (loading) {
    return <div>Loading...</div>;
  }

  return (
    <div className="p-4">
      <h1 className="text-2xl font-bold mb-4">Email Management System</h1>

      {errorMessages.general && (
        <div className="text-red-500 mb-4">{errorMessages.general}</div>
      )}

      <div className="space-y-4">
        {excelData.map((row) => (
          <div key={row.Sno} className="border p-4 rounded">
            <div className="grid grid-cols-2 gap-2">
              <div>SNO: {row.Sno}</div>
              <div>Name: {row.Name}</div>
              <div>Email: {row.Email}</div>
              <div>Company: {row.Company}</div>
              <div>Position: {row.ApplyingPosition}</div>
              <div>Additional Info: {row.AdditionalInfo}</div>
              <div>Reason: {row.ReasonForContact}</div>
            </div>

            <div className="mt-2 space-x-2">
              <button
                onClick={() => handleGenerateMail(row)}
                disabled={processingRows.has(row.Sno)}
                className="bg-blue-500 text-white px-4 py-2 rounded disabled:opacity-50"
              >
                {processingRows.has(row.Sno)
                  ? "Generating..."
                  : "Generate Mail"}
              </button>
              <button
                onClick={() => handleSendMail(row)}
                disabled={
                  processingRows.has(row.Sno) || !generatedEmails[row.Sno]
                }
                className="bg-green-500 text-white px-4 py-2 rounded disabled:opacity-50"
              >
                {processingRows.has(row.Sno) ? "Sending..." : "Send Mail"}
              </button>
            </div>

            {errorMessages[row.Sno] && (
              <div className="mt-2 text-red-500">{errorMessages[row.Sno]}</div>
            )}

            {successMessages[row.Sno] && (
              <div className="mt-2 text-green-500">
                {successMessages[row.Sno]}
              </div>
            )}

            {generatedEmails[row.Sno] && (
              <div className="mt-2">
                <div>
                  <label className="block">Subject:</label>
                  <input
                    type="text"
                    value={generatedEmails[row.Sno].Subject}
                    onChange={(e) =>
                      handleEmailContentChange(
                        row.Sno,
                        "Subject",
                        e.target.value
                      )
                    }
                    className="w-full p-2 border rounded"
                  />
                </div>
                <div className="mt-2">
                  <label className="block">Body:</label>
                  <textarea
                    value={generatedEmails[row.Sno].HTML}
                    onChange={(e) =>
                      handleEmailContentChange(row.Sno, "HTML", e.target.value)
                    }
                    className="w-full p-2 border rounded"
                    rows={4}
                  />
                </div>
              </div>
            )}
          </div>
        ))}
      </div>
    </div>
  );
};

export default App;
