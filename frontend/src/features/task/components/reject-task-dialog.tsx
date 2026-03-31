import { useState } from 'react';
import { X } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog';
import { Label } from '@/components/ui/label';

export function RejectTaskDialog({ onReject }: { onReject: (reason: string) => void }) {
  const [reason, setReason] = useState('');

  const handleReject = () => {
    onReject(reason);
    setReason('');
  };

  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button size="sm" variant="destructive">
          <X className="mr-2 h-4 w-4" /> Reprovar
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Reprovar tarefa</DialogTitle>
          <DialogDescription>Explique para o responsável o que precisa ser ajustado antes de uma nova revisão.</DialogDescription>
        </DialogHeader>
        <div className="space-y-4 pt-4">
          <div className="space-y-2">
            <Label htmlFor="task-reject-reason">Motivo da reprovação</Label>
            <textarea
              id="task-reject-reason"
              value={reason}
              onChange={(event) => setReason(event.target.value)}
              placeholder="Ex.: faltou validar o fluxo quando a task volta para To Do."
              className="min-h-28 w-full rounded-md border border-slate-200 px-3 py-2 text-sm outline-none transition focus:border-slate-400"
            />
          </div>
          <Button className="w-full" variant="destructive" onClick={handleReject} disabled={!reason.trim()}>
            Confirmar reprovação
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
