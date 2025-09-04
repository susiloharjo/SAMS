'use client';

import { useState, useRef, useEffect } from 'react';
import { Send, Bot, User, Sparkles, MessageCircle } from 'lucide-react';
import { api } from '@/utils/api';

interface Message {
  role: 'user' | 'assistant';
  content: string;
}

export default function AIChatPage() {
  const [messages, setMessages] = useState<Message[]>([
    {
      role: 'assistant',
      content: `Hello! I'm SAMS your AI Assistant. I can help you with asset management questions and provide real-time information about your inventory.

**What I can help you with:**
• **Asset Information**: Search for specific assets by name, model, or serial number
• **Asset Categories**: Find assets by type (IT Equipment, Vehicles, Tools)
• **Asset Status**: Find active, maintenance, or disposed assets  
• **Asset Location**: Find assets by location, building, or room
• **Asset Summaries**: Get total counts, values, and overviews

**Try asking questions like:**
• "What's the total value of IT equipment?"
• "How many vehicles do we have?"
• "Show me all active assets"
• "Tell me about Samsung Galaxy Tab S7"

How can I assist you today?`
    }
  ]);
  const [input, setInput] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!input.trim() || isLoading) return;

    const userMessage: Message = { role: 'user', content: input };
    setMessages((prev) => [...prev, userMessage]);
    setInput('');
    setIsLoading(true);

    try {
      const response = await api.post('/api/v1/ai/query', { message: input });
      const data = await response.json();
      
      const assistantMessage: Message = {
        role: 'assistant',
        content: data.response || 'Sorry, I could not get a response.',
      };
      setMessages((prev) => [...prev, assistantMessage]);

    } catch (error) {
      const errorMessage: Message = {
        role: 'assistant',
        content: error instanceof Error ? error.message : 'An unknown error occurred.',
      };
      setMessages((prev) => [...prev, errorMessage]);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="flex flex-col h-[calc(100vh-4rem)] bg-gradient-to-br from-slate-50 to-blue-50">
      {/* Modern Header */}
      <header className="p-6 border-b bg-white/80 backdrop-blur-sm border-gray-200/60 shadow-sm">
        <div className="flex items-center gap-4">
          <div className="w-12 h-12 bg-gradient-to-br from-blue-600 to-purple-600 rounded-2xl flex items-center justify-center shadow-lg">
            <Sparkles className="w-6 h-6 text-white" />
          </div>
          <div>
            <h1 className="text-2xl font-bold bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent">
              SAMS AI Assistant
            </h1>
            <p className="text-gray-600 mt-1">Powered by Gemini AI • Your intelligent asset management companion</p>
          </div>
        </div>
      </header>

      {/* Chat Messages */}
      <div className="flex-1 overflow-y-auto p-6 space-y-6">
        {messages.map((msg, index) => (
          <div key={index} className={`flex items-start gap-4 ${msg.role === 'user' ? 'justify-end' : 'justify-start'}`}>
            {msg.role === 'assistant' && (
              <div className="w-10 h-10 bg-gradient-to-br from-blue-500 to-blue-600 rounded-full flex items-center justify-center text-white flex-shrink-0 shadow-lg">
                <Bot className="w-5 h-5" />
              </div>
            )}
            
            <div className={`max-w-2xl px-6 py-4 rounded-2xl shadow-sm ${
              msg.role === 'user'
                ? 'bg-gradient-to-r from-blue-600 to-blue-700 text-white ml-12'
                : 'bg-white border border-gray-200 shadow-md'
            }`}>
              <p className="text-sm leading-relaxed whitespace-pre-wrap">{msg.content}</p>
            </div>
            
            {msg.role === 'user' && (
              <div className="w-10 h-10 bg-gradient-to-br from-gray-600 to-gray-700 rounded-full flex items-center justify-center text-white flex-shrink-0 shadow-lg">
                <User className="w-5 h-5" />
              </div>
            )}
          </div>
        ))}
        
        {/* Loading Indicator */}
        {isLoading && (
          <div className="flex items-start gap-4">
            <div className="w-10 h-10 bg-gradient-to-br from-blue-500 to-blue-600 rounded-full flex items-center justify-center text-white flex-shrink-0 shadow-lg">
              <Bot className="w-5 h-5" />
            </div>
            <div className="bg-white border border-gray-200 rounded-2xl px-6 py-4 shadow-md">
              <div className="flex items-center gap-2">
                <div className="flex space-x-1">
                  <div className="w-2 h-2 bg-blue-500 rounded-full animate-bounce"></div>
                  <div className="w-2 h-2 bg-blue-500 rounded-full animate-bounce" style={{ animationDelay: '0.1s' }}></div>
                  <div className="w-2 h-2 bg-blue-500 rounded-full animate-bounce" style={{ animationDelay: '0.2s' }}></div>
                </div>
                <span className="text-gray-500 text-sm ml-2">AI is thinking...</span>
              </div>
            </div>
          </div>
        )}
        
        <div ref={messagesEndRef} />
      </div>

      {/* Modern Input Form */}
      <div className="p-6 bg-white/80 backdrop-blur-sm border-t border-gray-200/60">
        <form onSubmit={handleSubmit} className="max-w-4xl mx-auto">
          <div className="relative">
            <input
              type="text"
              value={input}
              onChange={(e) => setInput(e.target.value)}
              placeholder="Ask about asset values, categories, locations, or search for specific assets..."
              className="w-full px-6 py-4 pr-16 border border-gray-200 rounded-2xl focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent shadow-sm text-gray-700 placeholder-gray-400 transition-all duration-200 hover:shadow-md focus:shadow-lg"
              disabled={isLoading}
            />
            <button
              type="submit"
              disabled={isLoading || !input.trim()}
              className="absolute right-2 top-1/2 transform -translate-y-1/2 w-12 h-12 bg-gradient-to-r from-blue-600 to-purple-600 text-white rounded-xl hover:from-blue-700 hover:to-purple-700 disabled:from-gray-300 disabled:to-gray-400 flex items-center justify-center shadow-lg transition-all duration-200 hover:shadow-xl disabled:shadow-sm"
            >
              <Send className="w-5 h-5" />
            </button>
          </div>
          
          {/* Helpful Suggestions */}
          <div className="mt-4 flex flex-wrap gap-2">
            <span className="text-xs text-gray-500 mr-2">Try asking:</span>
            {[
              "What's the total value of IT equipment?",
              "How many vehicles do we have?",
              "Show me all active assets",
              "Find assets in Jakarta",
              "What's our total asset value?",
              "Show me maintenance assets",
              "Tell me about Samsung Galaxy Tab S7",
              "Assets by department"
            ].map((suggestion, index) => (
              <button
                key={index}
                type="button"
                onClick={() => setInput(suggestion)}
                className="text-xs px-3 py-1 bg-gray-100 hover:bg-gray-200 text-gray-600 rounded-full transition-colors duration-200 cursor-pointer"
              >
                {suggestion}
              </button>
            ))}
          </div>
        </form>
      </div>
    </div>
  );
}
